package golang_checker

import (
	"fmt"
	"github.com/goplume/core/health"
	"github.com/goplume/core/utils/logger"
	"runtime"
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

	this.Log.RLog.Info("Check golang ..")

	health := health.NewHealth()

	golangVersion := fmt.Sprintf("%s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	health.Up()
	this.Log.RLog.Info(golangVersion)
	health.AddInfo("version", golangVersion)

	return health
}
