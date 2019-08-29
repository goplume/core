package kafka

type KafkaConfig struct {
    ConfigurationName string
    Enable            bool
    ClientID          string
    Brokers           string
    Version           string
    Topic             string
    Verbose           bool
    Partition         int32
    ProducerEnable    bool
    ConsumerEnable    bool
    Consumer          ConsumerConfiguration
    ActionName        string
}

type ConsumerConfiguration struct {
    Group             string
    Oldest            bool
    AcceptMessageFunc AcceptMessageFunc
}

var SampleKafkaConfig = KafkaConfig{
    Brokers: "127.0.0.1:9092",
    Version: "",
    Topic:   "sarama",
    Verbose: true,
    Consumer: ConsumerConfiguration{
        Group:  "example",
        Oldest: true,
    },
}
