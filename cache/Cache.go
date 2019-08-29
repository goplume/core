package cache

import (
	"encoding/json"
	"fmt"
	"github.com/goplume/core/fault"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"time"
)

type Cache struct {
	//Log                    *logger.Logger `inject:""`
	RedisClient            *redis.Client `inject:""`
	Enabled                bool          `inject:""`
	EntityCacheLifetimeSec int           `inject:""`
	KeyNamePrefix          string
}

func (this *Cache) Get(
	rlog logrus.FieldLogger,
	entityRef string, entity interface{},
) (
	bool, fault.TypedError,
) {
	found, err := this.GetEntity(rlog, entityRef, entity)
	return found, err
}

func (this *Cache) Put(
	rlog logrus.FieldLogger,
	entityRef string,
	entity interface{},
) {
	this.PutEntity(rlog, entityRef, entity)
}

func (this *Cache) Evict(
	rlog logrus.FieldLogger,
	entityRef string,
) {
	this.EvictEntity(rlog, entityRef)
}

func (this *Cache) GetEntity(
	rlog logrus.FieldLogger,
	entityRef string,
	entity interface{},
) (
	bool,
	fault.TypedError,
) {

	key := this.keyName(entityRef)
	rlog.Info(fmt.Sprintf("CACHE: Read %s", key))
	serialized, redisError := this.RedisClient.Get(key).Result()

	if redisError == redis.Nil {
		rlog.Info("CACHE: Not found " + key)
		// not found
		return false, nil
	} else if redisError != nil {
		return false, fault.ExceptionInternalError("Redis client error: " + redisError.Error())
	} else {
		redisError = json.Unmarshal([]byte(serialized), entity)
		if redisError != nil {
			return false, fault.ExceptionInternalError("Redis client error: " + redisError.Error())
		}
		rlog.Info(fmt.Sprintf("CACHE: Get %s", key))
		return true, nil
	}
}

func (this *Cache) PutEntity(
	rlog logrus.FieldLogger,
	entityRef string,
	entity interface{},
) fault.TypedError {
	if this.Enabled == false {
		return nil
	}

	serialized, marshalError := json.Marshal(entity)
	if marshalError != nil {
		return fault.ExceptionInternalError("Redis client error: " + marshalError.Error())
	}
	key := this.keyName(entityRef)
	duration := time.Duration(this.EntityCacheLifetimeSec) * time.Second
	rlog.Info(fmt.Sprintf("CACHE: Put %s on %s", key, duration))
	set := this.RedisClient.Set(
		key, serialized, duration)
	setRedisErr := set.Err()
	if setRedisErr != nil {
		return fault.ExceptionInternalError("Redis client error: " + setRedisErr.Error())
	}

	this.RedisClient.Expire(key, duration)
	//ttl := this.RedisClient.TTL(key)
	//this.RedisClient.Persist(key)

	rlog.Info(fmt.Sprintf("CACHE: Put %s Duration: %s", key, duration))
	return nil
}

func (this *Cache) EvictEntity(
	rlog logrus.FieldLogger,
	entityRef string,
) fault.TypedError {
	key := this.keyName(entityRef)
	del := this.RedisClient.Del(key)
	//i, err := del.Result()
	err := del.Err()
	if err != nil {
		return fault.ExceptionInternalError("Redis client error: " + err.Error())
	}
	rlog.Info("CACHE: Evict " + key)
	return nil
}

func (this *Cache) keyName(entityRef string) string {
	return this.KeyNamePrefix + "@" + entityRef;
}
