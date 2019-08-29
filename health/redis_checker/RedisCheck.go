package redis_checker

import (
	"github.com/goplume/core/health"
	"github.com/go-redis/redis"
	"strings"
)

// Redis is a interface used to abstract the access of the Version string
type Redis interface {
	GetVersion() (string, error)
}

// Checker is a checker that check a given redis
type Checker struct {
	RedisClient *redis.Client
}

// NewChecker returns a new redis.Checker
func NewChecker(network, addr string) Checker {

	// todo restore
	//return Checker{Redis: NewRedigo(network, addr)}
	return Checker{}
}

// NewCheckerWithRedis returns a new redis.Checker configured with a custom Redis implementation
func NewCheckerWithRedis(redisClient *redis.Client) Checker {
	return Checker{RedisClient: redisClient}
}

// Check obtain the version string from redis info command
func (this Checker) Check() health.Health {
	healthData := health.NewHealth()

	healthData.Down()
	if this.RedisClient == nil {
		healthData.AddInfo("error", "Redis is not configuration")
	}

	redisAddress := this.RedisClient.String()
	if len(redisAddress) > 0 {
		spaceIndex := strings.Index(redisAddress, " ")
		redisAddress = redisAddress[6:spaceIndex]
		health.TelnetCheck("", "redis://"+redisAddress, &healthData)
	}

	pong, err := this.RedisClient.Ping().Result()
	//version, err := this.Redis.Ping()

	if err != nil {
		healthData.Down().AddInfo("error", err.Error())
	} else {
		healthData.Up()
	}

	if pong != "" {
		healthData.AddInfo("state", pong)
	}

	//healthData.Up().AddInfo("version", version)

	return healthData
}
