package config_checker

import (
	"github.com/goplume/core/configuration"
	"github.com/goplume/core/health"
	"github.com/goplume/core/utils/logger"
)

type Checker struct {
	Log *logger.Logger
}

func NewChecker(log *logger.Logger) Checker {
	return Checker{
		Log: log,
	}
}

func (this Checker) Check() health.Health {

	this.Log.RLog.Info("Check configuration ..")

	health := health.NewHealth()

	health.Up()
	if len(configuration.LoadConfigurationError) > 0 {
		health.Down()
		health.AddInfo("load-configuration-error", configuration.LoadConfigurationError)
	}

	return health
}
