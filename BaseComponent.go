package core

import (
    "github.com/goplume/core/configuration"
    "github.com/goplume/core/factory/factory_cache"
    "github.com/goplume/core/factory/factory_datasource"
    "github.com/goplume/core/factory/factory_gin"
    "github.com/goplume/core/factory/factory_gormsource"
    "github.com/goplume/core/factory/factory_kafka"
    "github.com/goplume/core/factory/factory_logger"
    "github.com/goplume/core/factory/factory_oauth"
    "github.com/goplume/core/factory/factory_redirect"
    "github.com/goplume/core/factory/factory_redis"
    "github.com/goplume/core/factory/factory_restclient"
    "github.com/goplume/core/factory/factory_swagger"
    "github.com/goplume/core/health"
    "github.com/goplume/core/rest_api"
    "github.com/goplume/core/utils/logger"
    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/swaggo/gin-swagger/swaggerFiles"
    "net/http"
)

type AbstractComponent struct {
    Log *logger.Logger
}

type BaseComponent struct {
    AbstractComponent
    Address           string
    FactoryGinRouter  factory_gin.FactoryGinRouter
    FactoryRestClient factory_restclient.FactoryRestClient
    FactoryLogger     factory_logger.FactoryLogger
    FactoryRedis      factory_redis.FactoryRedis
    FactoryCache      factory_cache.FactoryCache
    FactoryOAuth      factory_oauth.FactoryOAuth
    FactoryGormSource factory_gormsource.FactoryGormSource
    FactorySwagger    factory_swagger.FactorySwagger
    FactoryKafka      factory_kafka.FactoryKafka
    FactoryDataSource factory_datasource.FactoryDataSource
    FactoryRedirect   factory_redirect.FactoryRedirect
    ApiVersion        string
    swaggerInfo       *SI
}

type SI struct {
    Version     string
    Host        string
    BasePath    string
    Schemes     []string
    Title       string
    Description string
}

func (this *BaseComponent) InitBaseComponent(
    serviceName string,
    version string,
// deprecated
    swaggerInfo *SI,
) {
    this.ApiVersion = version
    this.Log = this.FactoryLogger.CreateLogger(serviceName /*, elasticRestClient*/)
    rlog := this.Log.RLog

    rlog.Info("Start '" + serviceName + "'; Version: " + version + " ...")

    this.FactoryRedis.Log = this.Log
    this.FactoryCache.Log = this.Log
    this.FactoryOAuth.Log = this.Log
    this.FactoryKafka.Log = this.Log
    this.FactorySwagger.Log = this.Log
    this.FactoryRedirect.Log = this.Log
    this.FactoryGinRouter.Log = this.Log
    this.FactoryGormSource.Log = this.Log
    this.FactoryRestClient.Log = this.Log
    this.FactoryDataSource.Log = this.Log
}

func (this *BaseComponent) PostConstructComponent(
    serviceName string,
    swaggerInfo *SI,
    router Router,
    healthFunction health.HealthFunc,
) {
    config := configuration.NewServiceConfiguration(serviceName, "%s", this.Log)
    this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)
    this.swaggerInfo = swaggerInfo

    if config.GetBool("swagger-info.enable") {
        swagInfo := this.FactorySwagger.CreateSwagger(serviceName, this.ApiVersion)

        // todo ???
        swaggerInfo.Title = swagInfo.Title
        swaggerInfo.Description = swagInfo.Description
        swaggerInfo.Version = swagInfo.Version
        swaggerInfo.Host = swagInfo.Host
        swaggerInfo.BasePath = swagInfo.BasePath

        router.Get(rest_api.Endpoint_Log, BrowseLog)
        router.Get(rest_api.Endpoint_Api, this.RedirectToApi)
        //router.Get(rest_api.Endpoint_Docs, this.DocsRender)
        //router.Get(rest_api.Endpoint_Api_Slash, this.RedirectToApi)
        router.Get(rest_api.Endpoint_Root, this.RedirectToApi)
        router.Get(rest_api.Endpoint_Api_Any, ginSwagger.WrapHandler(swaggerFiles.Handler))
    }

    if config.GetBool("health.enable") {
        if healthFunction != nil {
            router.Get(rest_api.Endpoint_Health, HealthHandler(healthFunction))
        }
    }

    stateFilter := rest_api.StateFilter{}
    DoCurrentSatet := stateFilter.DoCurrentSatet
    router.AddFilter(DoCurrentSatet)

    if config.GetBool("auth.enable") {
        restOauthClient := this.FactoryRestClient.CreateRestClient(serviceName, "oauth")
        oauthClient := this.FactoryOAuth.CreateOAuth(serviceName, restOauthClient.HttpClient)
        authFilter := &rest_api.AuthFilter{
            OAuthClient:           oauthClient,
            Log:                   this.Log,
            Scope:                 config.GetString("auth.scope"),
            ExcludePath:           config.GetStrings("auth.excludePaths"),
            SkipCheckMpanPath:     config.GetStrings("auth.uncheckMpan"),
            PermanentAccessTokens: config.GetMap("auth.accessTokens"),
        }
        if config.GetBool("auth.mock.enable") {
            authFilterMock := &rest_api.AuthFilterMock{
                Mpan: config.GetString("auth.mock.mpan"),
            }
            authFilter.Mock = authFilterMock
        }

        router.AddFilter(authFilter.DoAuth)
    }

    // Setting api version for response
    router.AddFilter(func(context *gin.Context) {
        context.Set("api_version", this.ApiVersion)
    })

    router.GetRouterGroup().StaticFS(rest_api.Endpoint_Static, gin.Dir("static", true))

    router.InitRouter()
    this.Log.RLog.Info("Started")
}

func (this *BaseComponent) RedirectToApi(ctx *gin.Context) {
    ctx.Redirect(http.StatusPermanentRedirect, this.swaggerInfo.BasePath+rest_api.Endpoint_Api_Index)
}
