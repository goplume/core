package factory_logger

import (
	"github.com/goplume/core/configuration"
	"github.com/goplume/core/utils/logger"
	"github.com/sirupsen/logrus"
)

type FactoryLogger struct {
}

func (this *FactoryLogger) InitFactory() {
}

func (this *FactoryLogger) CreateLogger(
	serviceName string,
//elasticRestClient *rest_client.MerchantServiceRestClient,
) *logger.Logger {

	config := configuration.NewServiceConfiguration(serviceName, "logger.%s", nil)

	url_log := config.GetString("url")
	//rlog.Info("url_log := " + url_log)

	strategyName := config.GetString("strategy")
	if strategyName == "" {
		strategyName = string(logger.CONSOLE_STRATEGY)
	}

	var strategy logger.LoggingStrategy

	switch strategyName {
	case string(logger.ELK_STRATEGY):
		strategy = logger.ELK_STRATEGY
	case string(logger.CONSOLE_STRATEGY):
		strategy = logger.CONSOLE_STRATEGY
	case string(logger.FILE_STRATEGY):
		strategy = logger.FILE_STRATEGY
	case string(logger.FILE_ELK_STRATEGY):
		strategy = logger.FILE_ELK_STRATEGY
	}

	var nlog = logrus.New()
	nlog.SetLevel(logrus.DebugLevel)
	switch strategy {
	case logger.ELK_STRATEGY:
		logger.SetupLogrusToFile(nlog, url_log)
	case logger.CONSOLE_STRATEGY:
		logger.SetupLogrusToConsole(nlog)
	case logger.FILE_STRATEGY:
		logger.SetupLogrusToFile(nlog, url_log)
	case logger.FILE_ELK_STRATEGY:
		logger.SetupLogrusToFile(nlog, url_log)
	default:
		logger.SetupLogrusToFile(nlog, url_log)
	}
	rlog := nlog.WithField("component", serviceName)

	fields := config.GetMap("fields")
	if fields != nil {
		for k, v := range fields {
			rlog = rlog.WithField(k, v)
		}
	}

	rlog.Info("strategy := " + strategyName)
	rlog.Info("Read configuration from context " + config.ConfigurationContext)

	logger := &logger.Logger{
		RLog:       rlog,
		Strategy:   strategy,
		Component:  serviceName,
		ElasticUrl: url_log,
		//ElasticRestClient: elasticRestClient,
		//ElkLogger: logger.NewLogger(url_log),
	}

	return logger
}
