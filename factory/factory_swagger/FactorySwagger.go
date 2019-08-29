package factory_swagger

import (
    "github.com/goplume/core/configuration"
    "github.com/goplume/core/utils/logger"
)

type FactorySwagger struct {
    Log *logger.Logger
}

func (this *FactorySwagger) InitFactory() {
}

func (this *FactorySwagger) CreateSwagger(
    serviceName string,
    version string,
) (
    swagInfo struct {
    Version     string
    Host        string
    BasePath    string
    Title       string
    Description string
},
) {
    config := configuration.NewServiceConfiguration(serviceName, "swagger-info.%s", this.Log)
    this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

    swagInfo.Title = config.GetString("Title")
    swagInfo.Description = config.GetString("Description")
    swagInfo.Version = version
    swagInfo.Host = config.GetString("Host")
    swagInfo.BasePath = config.GetString("BasePath")

    return
}
