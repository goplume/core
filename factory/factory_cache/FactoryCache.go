package factory_cache

import (
	"github.com/goplume/core/cache"
	"github.com/goplume/core/configuration"
	"github.com/goplume/core/utils/logger"
	"github.com/go-redis/redis"
)

type FactoryCache struct {
	Log *logger.Logger
}

func (this *FactoryCache) InitFactory() {
}

func (this *FactoryCache) CreateCache(
	serviceName string,
	cacheName string,
	RedisClient *redis.Client,
) *cache.Cache {
	config := configuration.NewServiceConfiguration(serviceName, "cache."+cacheName+".%s", this.Log)
	this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

	return &cache.Cache{
		RedisClient:            RedisClient,
		Enabled:                config.GetBool("enable"),
		EntityCacheLifetimeSec: config.GetInt("lifetime-sec"),
		KeyNamePrefix:          config.GetString("key"),
	}

}
