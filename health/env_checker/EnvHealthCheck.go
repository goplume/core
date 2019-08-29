package env_checker

import (
	"fmt"
	"github.com/goplume/core/health"
	utils2 "github.com/goplume/core/utils"
	"github.com/goplume/core/utils/logger"
	"github.com/goplume/core/utils/third/goInfo"
	"os"
	"strings"
	"time"
)

var start_time = time.Now()

type Checker struct {
	Log *logger.Logger
}

func NewChecker(log *logger.Logger) Checker {
	return Checker{
		Log: log,
	}
}

func (this Checker) Check() health.Health {

	this.Log.RLog.Info("Check enviroments ..")

	health := health.NewHealth()

	health.Up()
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		envInfo := fmt.Sprint(variable[0], "=>", variable[1])
		this.Log.RLog.Info(envInfo)
		health.AddInfo(variable[0], variable[1])
	}

	gi := goInfo.GetInfo()

	health.AddInfo("start-time", start_time.Format(utils2.DATETIME_FORMAT_LAYOUT))
	health.AddInfo("execution-time", fmt.Sprintf("%v", time.Since(start_time)))
	health.AddInfo("os.GoOS", gi.GoOS)
	health.AddInfo("os.Kernel", gi.Kernel)
	health.AddInfo("os.Core", gi.Core)
	health.AddInfo("os.Platform", gi.Platform)
	health.AddInfo("os.OS", gi.OS)
	health.AddInfo("os.Hostname", gi.Hostname)
	health.AddInfo("os.CPUs", gi.CPUs)

	return health
}
