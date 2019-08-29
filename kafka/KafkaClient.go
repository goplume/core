package kafka

import (
	"github.com/goplume/core/fault"
	"github.com/goplume/core/utils/logger"
	"github.com/sirupsen/logrus"
)

type KafkaClient interface {
	Init(*logger.Logger, *KafkaConfig)
	KafkaConfig() *KafkaConfig
	Loop() fault.TypedError
	GroupLoop() fault.TypedError
	ConfigurationName() string
	Version() string
	Push([]byte) fault.TypedError
	PushText(string) fault.TypedError
	PushEvent(interface{}) fault.TypedError
	Ping() fault.TypedError
	Status() string
	IsWorked() bool
	GetMetrics() map[string]*Metrics
	GetHost() []string
	SetAcceptMessageFunc(AcceptMessageFunc)
}

type Metrics struct {
	ProduceMsgCount         uint64
	ConsumeMsgCount         uint64
	ProduceFailedMsgCount   uint64
	ProduceSuccessMsgCount  uint64
	HealthCheckFailedCount  uint64
	ConsumeFailedMsgCount   uint64
	ConsumeSuccessMsgCount  uint64
	HealthCheckCount        uint64
	HealthCheckSuccessCount uint64
}

type AcceptMessageFunc func(logrus.FieldLogger, []byte) (fault.TypedError, string)
