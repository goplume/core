package factory_redirect

import (
	"github.com/goplume/core/configuration"
	"github.com/goplume/core/utils/logger"
)

type FactoryRedirect struct {
	Log *logger.Logger
}

func (this FactoryRedirect) InitFactory() {
}

func (this FactoryRedirect) CreateRedirectPrefix(
	serviceName string,
) string {
	config := configuration.NewServiceConfiguration(serviceName, "redirect.%s", this.Log)

	this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

	redirectPrefixUrl := config.GetString("prefix-url")
	return redirectPrefixUrl
}
