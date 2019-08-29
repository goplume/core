package factory_kafka

import (
    "github.com/goplume/core/configuration"
    "github.com/goplume/core/kafka"
    "github.com/goplume/core/utils/logger"
)

type FactoryKafka struct {
    Log *logger.Logger
}

func (this *FactoryKafka) InitFactory() {
}

func (this *FactoryKafka) CreateProduceFactoryKafka(
    serviceName string,
    configurationName string,
    actionName string,
) (kafkaClient kafka.KafkaClient) {
    return this.CreateFactoryKafka(serviceName, configurationName, nil, actionName)
}

func (this *FactoryKafka) CreateFactoryKafka(
    serviceName string,
    configurationName string,
    acceptMessageFunc kafka.AcceptMessageFunc,
    actionName string,
) (kafkaClient kafka.KafkaClient) {
    config := configuration.NewServiceConfiguration(serviceName, "kafka."+configurationName+".%s", this.Log)
    this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

    kafkaConfig := &kafka.KafkaConfig{
        ConfigurationName: configurationName,
        Enable:            config.GetBoolD("Enable", false),
        ClientID:          config.GetString("ClientID"),
        Brokers:           config.GetString("Brokers"),
        Version:           config.GetString("Version"),
        Topic:             config.GetString("Topic"),
        Partition:         config.GetInt32("Partition"),
        ProducerEnable:    config.GetBool("Producer.enable"),
        ConsumerEnable:    config.GetBool("Consumer.enable"),
        Consumer: kafka.ConsumerConfiguration{
            Oldest:            config.GetBool("Consumer.Oldest"),
            Group:             config.GetString("Consumer.Group"),
            AcceptMessageFunc: acceptMessageFunc,
        },
        ActionName: actionName,
    }

    //impl := kafka.KafkaClientImpl{}
    kafkaClientImpl := &kafka.KafkaClientImpl{}
    if kafkaConfig.Enable {
        kafkaClientImpl.Init(
            this.Log,
            kafkaConfig,
        )
    }

    return kafkaClientImpl

}
