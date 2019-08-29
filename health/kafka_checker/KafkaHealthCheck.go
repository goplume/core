package kafka_checker

import (
	"github.com/goplume/core/health"
	"github.com/goplume/core/kafka"
	"github.com/goplume/core/utils/logger"
)

type Checker struct {
	KafkaClient kafka.KafkaClient
	Log         *logger.Logger
}

func NewChecker(
	log *logger.Logger,
	KafkaClient kafka.KafkaClient,
) Checker {
	return Checker{
		KafkaClient: KafkaClient,
		Log:         log,
	}
}

func (this Checker) Check() health.Health {
	healthData := health.NewHealth()
	healthData.Up()

	if this.KafkaClient == nil || this.KafkaClient.ConfigurationName() == "" {
		healthData.AddInfo("status", "not configured")
		healthData.Up()
		return healthData
	}

	healthData.AddInfo("configuration-name", this.KafkaClient.ConfigurationName())
	healthData.AddInfo("topic", this.KafkaClient.KafkaConfig().Topic)
	healthData.AddInfo("enable", this.KafkaClient.KafkaConfig().Enable)
	healthData.AddInfo("ConsumerEnable", this.KafkaClient.KafkaConfig().ConsumerEnable)
	healthData.AddInfo("ProducerEnable", this.KafkaClient.KafkaConfig().ProducerEnable)
	//healthData.AddInfo("config", this.KafkaClient.KafkaConfig())
	healthData.AddInfo("version", this.KafkaClient.Version())
	healthData.AddInfo("status", this.KafkaClient.Status())
	healthData.AddInfo("metrics", this.KafkaClient.GetMetrics())
	this.KafkaClient.Ping()
	if !this.KafkaClient.IsWorked() {
		healthData.Down()
	}

	if (this.KafkaClient.GetHost() != nil) && (len(this.KafkaClient.GetHost()) > 0) {
		health.TelnetCheck("", "kafka://"+this.KafkaClient.GetHost()[0], &healthData)
	}

	return healthData
}
