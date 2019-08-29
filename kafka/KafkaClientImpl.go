package kafka

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "encoding/json"
    "fmt"
    "github.com/goplume/core/fault"
    "github.com/goplume/core/utils/logger"
    "github.com/Shopify/sarama"
    "io/ioutil"
    "strings"
    "time"
)

var KafkaClientRef = &KafkaClientImpl{}

type KafkaClientImpl struct {
    kafkaConfig       *KafkaConfig
    Log               *logger.Logger
    consumer          sarama.Consumer
    consumerGroup     sarama.ConsumerGroup
    producer          sarama.AsyncProducer
    client            sarama.Client
    acceptMessageFunc AcceptMessageFunc
    groupHandler      *ConsumG
    config            *sarama.Config
    Topic             string
    Partition         int32
    Offset            string
    MsgCount          int64
    IsWorkedStatus    bool
    StatusError       error
    kafkaVersion      string
    Metrics           map[string]*Metrics
    Brokers           []string
    //groupHandler  *sarama.ConsumerGroupHandler
}

func (this *KafkaClientImpl) Init(
    log *logger.Logger,
    kafkaConfig *KafkaConfig,
) {
    var err error
    this.Metrics = make(map[string]*Metrics)
    this.kafkaConfig = kafkaConfig
    this.IsWorkedStatus = true
    this.Topic = kafkaConfig.Topic
    this.Partition = kafkaConfig.Partition
    this.Log = log
    rlog := this.Log.RLog.WithField("action", kafkaConfig.ActionName)

    if !(kafkaConfig.ProducerEnable || kafkaConfig.ConsumerEnable) {
        err = fault.ExceptionInternalError(
            "Kafka Client: Must configure producer or consumer or both")
        this.StatusError = err
        this.IsWorkedStatus = false
        rlog.Error(err)
        return
    }

    /**
     * Construct a new Sarama configuration.
     * The Kafka cluster version has to be defined before the consumer/producer is initialized.
     */
    this.kafkaVersion = kafkaConfig.Version
    kafkaVersion, err := sarama.ParseKafkaVersion(kafkaConfig.Version)
    if err != nil {
        err = fault.ExceptionInternalError(
            "Kafka Client: Invalid kafka version " + kafkaConfig.Version +
                ". Error: " + err.Error())
        this.StatusError = err
        this.IsWorkedStatus = false
        rlog.Error(err)
        return
    }

    this.config = sarama.NewConfig()
    this.config.ClientID = kafkaConfig.ClientID
    //this.config.
    this.config.Net.ReadTimeout = 5 * time.Second
    this.config.Net.DialTimeout = 5 * time.Second
    this.config.Net.WriteTimeout = 5 * time.Second

    // todo add logger
    //sarama.Logger = this.Log

    this.config.Version = kafkaVersion

    this.Brokers = strings.Split(kafkaConfig.Brokers, ",", )
    this.client, err = sarama.NewClient(strings.Split(kafkaConfig.Brokers, ",", ), this.config)
    if err != nil {
        err = fault.ExceptionInternalError("Kafka Client: " + err.Error())
        this.StatusError = err
        this.IsWorkedStatus = false
        rlog.Error(err)
        return
    }

    if kafkaConfig.ProducerEnable {
        this.producer, err = sarama.NewAsyncProducerFromClient(this.client)
        if err != nil {
            err = fault.ExceptionInternalError("Kafka Client: " + err.Error())
            this.StatusError = err
            this.IsWorkedStatus = false
            rlog.Error(err)
            return
        }
    }

    if kafkaConfig.ConsumerEnable {
        if kafkaConfig.Consumer.Oldest {
            this.config.Consumer.Offsets.Initial = sarama.OffsetOldest
        } else {
            this.config.Consumer.Offsets.Initial = sarama.OffsetNewest
        }
        this.acceptMessageFunc = kafkaConfig.Consumer.AcceptMessageFunc

        this.groupHandler = &ConsumG{
            Log:               this.Log,
            AcceptMessageFunc: this.acceptMessageFunc,
            KafkaClient:       this,
        }

        if len(kafkaConfig.Consumer.Group) == 0 {
            this.consumer, err = sarama.NewConsumerFromClient(this.client)
            if err != nil {
                err = fault.ExceptionInternalError("Kafka Client: " + err.Error())
                this.StatusError = err
                this.IsWorkedStatus = false
                rlog.Error(err)
                return
            }

            go this.Loop()
        } else {
            this.consumerGroup, err = sarama.NewConsumerGroupFromClient(
                kafkaConfig.Consumer.Group, this.client)
            if err != nil {
                err = fault.ExceptionInternalError("Kafka Client: " + err.Error())
                this.StatusError = err
                this.IsWorkedStatus = false
                rlog.Error(err)
                return
            }

            go this.GroupLoop()
        }
    }
}

func (this *KafkaClientImpl) SetAcceptMessageFunc(acceptMessageFunc AcceptMessageFunc) {
    this.acceptMessageFunc = acceptMessageFunc
}

func (this *KafkaClientImpl) Loop() fault.TypedError {

    partitionConsumer, err := this.consumer.ConsumePartition(
        this.Topic,
        this.Partition,
        sarama.OffsetNewest,
    )
    //<-this.consumer.ready // Await till the consumer has been set up
    if err != nil {
        this.Log.RLog.Error(err)
        return fault.ExceptionInternalErrorE(err)
    }
    defer partitionConsumer.Close()

    for {
        for message := range partitionConsumer.Messages() {
            logFields := map[string]interface{}{
                "offset":    message.Offset,
                "topic":     message.Topic,
                "Partition": message.Partition,
            }

            for k, v := range message.Headers {
                logFields[string(k)] = v
            }

            rlog := this.Log.RLog.WithFields(logFields)

            this.acceptMessageFunc(rlog, message.Value)
        }
    }
    //for {
    //    msg := <-partitionConsumer.Messages()
    //    log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)
    //}

}

func (this *KafkaClientImpl) GroupLoop() fault.TypedError {
    ctx := context.Background()
    topics := []string{this.Topic}
    this.groupHandler.ready = make(chan bool, 0)
    err := this.consumerGroup.Consume(
        ctx,
        topics,
        this.groupHandler,
    )
    if err != nil {
        err = fault.ExceptionInternalError("Kafka Client: " + err.Error())
        this.StatusError = err
        this.Log.RLog.Error(err)
        return fault.ExceptionInternalErrorE(err)
    }
    return nil
}

func (this *KafkaClientImpl) Version() (version string) {
    return this.kafkaVersion
}

func (this *KafkaClientImpl) IsWorked() bool {
    return this.IsWorkedStatus;
}

func (this *KafkaClientImpl) GetMetrics() map[string]*Metrics {
    return this.Metrics;
}

func (this *KafkaClientImpl) GetTopicMetrics(metricTopicIdx string) *Metrics {
    metrics, existMap := this.Metrics[metricTopicIdx]
    if existMap == false {
        metrics = &Metrics{}
        this.Metrics[metricTopicIdx] = metrics
    }

    return metrics;
}

func (this *KafkaClientImpl) Status() (status string) {
    status = "Failed"
    if this.IsWorkedStatus {
        status = "Up"
    }
    if this.StatusError != nil {
        status = status + ": " + this.StatusError.Error()
    }

    return status
}

func (this *KafkaClientImpl) GetHost() []string {
    return this.Brokers;
}

func (this *KafkaClientImpl) KafkaConfig() *KafkaConfig {
    return this.kafkaConfig
}

func (this *KafkaClientImpl) ConfigurationName() string {
    if this.kafkaConfig == nil {
        return ""
    }
    return this.kafkaConfig.ConfigurationName;
}

func (this *KafkaClientImpl) Pong() (fault.TypedError, string) {
    this.IsWorkedStatus = true
    //metrics := this.GetTopicMetrics(this.Topic)
    //metrics.HealthCheckSuccessCount++
    this.Log.RLog.Info("Kafka: Health Check: Accept message successfully")
    return nil, ""
}

func (this *KafkaClientImpl) Ping() fault.TypedError {
    if this.producer == nil {
        return fault.ExceptionInternalError("Producer for kafka is not initialized")
    }

    metrics := this.GetTopicMetrics(this.Topic)
    metrics.HealthCheckCount++
    this.producer.Input() <- &sarama.ProducerMessage{
        Topic: this.Topic,
        Key:   nil,
        //Offset: sarama.OffsetNewest,
        Offset: sarama.OffsetOldest,
        Value:  sarama.StringEncoder("pong"),
        Headers: []sarama.RecordHeader{
            {
                Key:   []byte("type"),
                Value: []byte("healthcheck"),
            },
        },
    }
    this.Log.RLog.Info("kafka ping")

    return nil
}

func (this *KafkaClientImpl) PushEvent(
    message interface{},
) fault.TypedError {
    bytes, err := json.Marshal(message)
    if err != nil {
        // todo ???
        return fault.ExceptionInternalErrorE(err)
    }
    err = this.Push(bytes)
    if err != nil {
        return fault.ExceptionInternalErrorE(err)
    }
    return nil
}

func (this *KafkaClientImpl) PushText(message string) fault.TypedError {
    return this.Push([]byte(message))
}

func (this *KafkaClientImpl) Push(
    message []byte,
) fault.TypedError {
    if this.producer == nil {
        return fault.ExceptionInternalError("Producer for kafka is not initialized")
    }
    if len(message) > 0 {
        this.Log.RLog.Info(fmt.Sprintf("Push in topic %s", this.Topic))

        metrics := this.GetTopicMetrics(this.Topic)

        metrics.ProduceMsgCount++
        this.producer.Input() <- &sarama.ProducerMessage{
            Topic: this.Topic,
            Key:   nil,
            Value: sarama.ByteEncoder(message),
            Headers: []sarama.RecordHeader{
                {
                    Key:   []byte("type"),
                    Value: []byte("data"),
                },
            },
        }
        //this.GetMetrics.ProduceFailedMsgCount++
        metrics.ProduceSuccessMsgCount++
        this.Log.RLog.Info(fmt.Sprintf("Produced message: [%s]\n", message))

    }
    return nil
}

func (this *KafkaClientImpl) genTLSConfig(
    clientcertfile string,
    clientkeyfile string,
    cacertfile string,
) (*tls.Config, fault.TypedError) {
    // load client cert
    clientcert, err := tls.LoadX509KeyPair(clientcertfile, clientkeyfile)
    if err != nil {
        return nil, fault.ExceptionInternalErrorE(err)
    }

    // load ca cert pool
    cacert, err := ioutil.ReadFile(cacertfile)
    if err != nil {
        return nil, fault.ExceptionInternalErrorE(err)
    }
    cacertpool := x509.NewCertPool()
    cacertpool.AppendCertsFromPEM(cacert)

    // generate tlcconfig
    tlsConfig := tls.Config{}
    tlsConfig.RootCAs = cacertpool
    tlsConfig.Certificates = []tls.Certificate{clientcert}
    tlsConfig.BuildNameToCertificate()
    // tlsConfig.InsecureSkipVerify = true // This can be used on test server if domain does not match cert:
    return &tlsConfig, nil
}

// Consumer represents a Sarama consumer group consumer
type ConsumG struct {
    Log               *logger.Logger
    ready             chan bool
    AcceptMessageFunc AcceptMessageFunc
    KafkaClient       *KafkaClientImpl
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ConsumG) Setup(sarama.ConsumerGroupSession) error {
    // Mark the consumer as ready
    close(consumer.ready)
    return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ConsumG) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (this *ConsumG) ConsumeClaim(
    session sarama.ConsumerGroupSession,
    claim sarama.ConsumerGroupClaim,
) error {

    // NOTE:
    // Do not move the code below to a goroutine.
    // The `ConsumeClaim` itself is called within a goroutine, see:
    for message := range claim.Messages() {
        metricTopicIdx := fmt.Sprintf("%s:%d", message.Topic, message.Partition)
        metrics := this.KafkaClient.GetTopicMetrics(metricTopicIdx)

        //message.Headers

        logFields := map[string]interface{}{
            "kafka_offset":    message.Offset,
            "kafka_topic":     message.Topic,
            "kafka_partition": message.Partition,
            "action":          this.KafkaClient.kafkaConfig.ActionName,
        }

        for _, v := range message.Headers {
            logFields["kafka_H["+string(v.Key)+"]"] = string(v.Value)
        }

        rlog := this.Log.RLog.WithFields(logFields)

        rlog.Info(fmt.Sprintf(
            "Kafka: Message claimed: value = %s, timestamp = %v, topic = %s",
            string(message.Value), message.Timestamp, message.Topic,
        ))

        var metadata string
        var err error
        if this.isHealthCheckMessage(message) {
            metrics.HealthCheckCount++
            err, metadata = this.KafkaClient.Pong()
            if err != nil {
                metrics.HealthCheckFailedCount++
                rlog.Error(err)
                return err
            }
            metrics.HealthCheckSuccessCount++
        } else {
            metrics.ConsumeMsgCount++
            err, metadata = this.AcceptMessageFunc(rlog, message.Value)
            if err != nil {
                metrics.ConsumeFailedMsgCount++
                rlog.Error(err)
                return err
            }
            metrics.ConsumeSuccessMsgCount++
        }

        { // final of message processing
            rlog.Info(fmt.Sprintf("Accept message with metadata size %d, data %v", message.Offset, metadata))
            session.MarkMessage(
                message,
                metadata, // todo необходимо как то обрабатывать метадату
            )
            session.MarkOffset(
                message.Topic,
                message.Partition,
                message.Offset,
                metadata,
            )

            rlog.Info(fmt.Sprintf("Commit Message with metadata size %d, data %v", message.Offset, metadata))
        }
    }

    return nil
}

func (this *ConsumG) isHealthCheckMessage(message *sarama.ConsumerMessage) bool {
    if len(message.Headers) >= 1 {
        for _, v := range message.Headers {
            if ("type" == string(v.Key) && "healthcheck" == string(v.Value)) {
                return true
            }
        }
    }
    return false
}
