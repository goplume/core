package factory_oauth

import (
    "github.com/goplume/core/cache"
    "github.com/goplume/core/configuration"
    "github.com/goplume/core/oaut_client"
    utils2 "github.com/goplume/core/utils"
    "github.com/goplume/core/utils/logger"
    "github.com/go-resty/resty/v2"
)

type FactoryOAuth struct {
    Log *logger.Logger
}

func (this FactoryOAuth) InitFactory() {
}

func (this FactoryOAuth) CreateOAuth(
    serviceName string,
    httpClient *resty.Client,
) (auth *oaut_client.OAuthClient) {
    return this.CreateCachedOAuth(serviceName, httpClient, nil)
}

func (this FactoryOAuth) CreateCachedOAuth(
    serviceName string,
    httpClient *resty.Client,
    cache *cache.Cache,
) (auth *oaut_client.OAuthClient) {
    config := configuration.NewServiceConfiguration(serviceName, "oauth.%s", this.Log)

    this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

    auth = &oaut_client.OAuthClient{
        Log:            this.Log,
        AuthEnabled:    config.GetBoolD("enable", true),
        GetTokenURL:    config.GetString("GetTokenURL"),
        CheckTokenURL:  config.GetString("CheckTokenURL"),
        RemoveTokenURL: config.GetString("RemoveTokenURL"),
        HttpClient:     httpClient,
        Cache:          cache,
    }

    scopes := config.GetMap("scopes")
    auth.Scopes = make(map[string]oaut_client.ClientSecret)
    for k, v := range scopes {
        c := v.(map[string]interface{})
        scope := oaut_client.ClientSecret{
            ClientID:     utils2.ToString(c["clientid"]),
            ClientSecret: utils2.ToString(c["clientsecret"]),
            Scope:        utils2.ToString(c["scope"]),
            GrantType:    utils2.ToString(c["granttype"]),
            Enable:       utils2.ToString(c["enable"]),
        }
        auth.Scopes[k] = scope
    }

    if auth.AuthEnabled == false {
        this.Log.RLog.Info("Authorization service disabled or not configured")
    }

    return
}
