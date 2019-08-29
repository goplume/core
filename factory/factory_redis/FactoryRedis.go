package factory_redis

import (
	"github.com/goplume/core/configuration"
	"github.com/goplume/core/utils/logger"
	"github.com/go-redis/redis"
)

type FactoryRedis struct {
	Log *logger.Logger
}

func (this *FactoryRedis) InitFactory() {
}

func (this FactoryRedis) CreateFactoryRedis(
	serviceName string,
) (redisClient *redis.Client) {
	config := configuration.NewServiceConfiguration(serviceName, "redis.%s", this.Log)
	this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

	RedisAddress := config.GetString("address")
	RedisPassword := config.GetString("password")
	RedisDb := config.GetInt("DB")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     RedisAddress,
		Password: RedisPassword,
		DB:       RedisDb,
	})

	// todo ???? restore
	//pong, err := redisClient.Ping().Result()
	//if err != nil {
	//    this.Log.RLog.Error(err)
	//}

	//this.Log.RLog.Info(pong)

	return redisClient
}
