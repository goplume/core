package rest_api

import (
    "github.com/goplume/core/health"
    "github.com/goplume/core/health/config_checker"
    "github.com/goplume/core/health/elk_checker"
    "github.com/goplume/core/health/env_checker"
    "github.com/goplume/core/health/golang_checker"
    "github.com/goplume/core/health/health_checker"
    "github.com/goplume/core/health/oauth_checker"
    "github.com/goplume/core/oaut_client"
    "github.com/goplume/core/rest_client"
    "github.com/goplume/core/utils/logger"
    "github.com/sirupsen/logrus"
)

type HealthApiService struct {
    RestClients map[string]*rest_client.RestClient
    OAuthClient *oaut_client.OAuthClient
    Log         *logger.Logger
    ServiceName string
    Versions    map[string]string
}

func (this *HealthApiService) AddRestClient(name string, restClient *rest_client.RestClient) {
    if this.RestClients == nil {
        this.RestClients = make(map[string]*rest_client.RestClient)
    }
	this.RestClients[name] = restClient
}

func (this *HealthApiService) Health() health.Health {
    compose := health.NewCompositeChecker()

    compose.AddChecker("ELK", elk_checker.NewChecker(this.Log))
    compose.AddChecker("Configuration", config_checker.NewChecker(this.Log))
    compose.AddChecker("golang", golang_checker.NewChecker(this.Log))
    compose.AddChecker("enviroment", env_checker.NewChecker(this.Log))
    compose.AddChecker("oauth", oauth_checker.NewChecker(this.Log, this.OAuthClient))
	for k, v := range this.RestClients {
		compose.AddChecker(k, health_checker.NewRestClientChecker(v))
	}

    compose.AddInfo("versions", this.Versions)

    h := compose.Check()

    this.Log.RLog.Info(h)
    logrus.Info(h)

    return h
}
