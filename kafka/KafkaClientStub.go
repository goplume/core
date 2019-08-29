// +build Stub

package kafka

import (
	"github.com/goplume/core/utils/logger"
)

var KafkaClientImpl = &KafkaClientStub{}

type KafkaClientStub struct {
	Log *logger.Logger
}

func (this *KafkaClientStub) Init(
	Log *logger.Logger,
	kafkaConfig *KafkaConfig,
) {
}

func (this *KafkaClientStub) Loop() (err error) {
	return nil
}

func (this *KafkaClientStub) GroupLoop() (err error) {
	return nil
}

func (this *KafkaClientStub) Version() (version string) {
	return ""
}

func (this *KafkaClientStub) IsWorked() bool {
	return false;
}

func (this *KafkaClientStub) GetMetrics() Metrics {
	return Metrics{};
}

func (this *KafkaClientStub) GetHost() []string {
	return []string{"localhost:9092"};
}

func (this *KafkaClientStub) Status() (version string) {
	return ""
}

func (this *KafkaClientStub) PushText(text string) error {
	return nil
}
