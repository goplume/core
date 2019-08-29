package elk_checker

import (
	"github.com/goplume/core/health"
	"github.com/goplume/core/utils/logger"
)

// Checker is a checker that check a given URL
type Checker struct {
	Log *logger.Logger
}

func NewChecker(log *logger.Logger) Checker {
	return Checker{
		Log: log,
	}
}

func (this Checker) Check() health.Health {

	this.Log.RLog.Info("Check ELK service ..")

	healthData := health.NewHealth()
	healthData.Up()

	health.TelnetCheck("", this.Log.ElasticUrl, &healthData)

	//if this.Log != nil && this.Log.ElasticRestClient != nil {
	//	this.Log.RLog.Info(" Test Logger")
	//	this.Log.RLog.Info("url" + this.Log.ElasticRestClient.HttpClient.HostURL)
	//
	//	//if err != nil {
	//	//	health.Down()
	//	//	health.AddInfo("err", err.Error())
	//	//	this.Log.RLog.Info("Checking failed " + err.Error())
	//	//} else {
	//	//	this.Log.RLog.Info("Check: ELK TransactionsService  UP")
	//	//	health.Up()
	//	//}
	//} else {
	//	health.AddInfo("status ","Elastic logger not configuration")
	//	health.Up()
	//}

	return healthData
}
