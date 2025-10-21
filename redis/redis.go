package redis

import (
	"os"
	"time"

	"franco/log"

	"github.com/go-redis/redis"
)

var Client *redis.Client

const Nil = redis.Nil

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}

func SetRedisString(key string, otp string, expiration time.Duration) error {
	err := Client.Set(key, otp, expiration).Err()
	if err != nil {
		log.LogDebug("Error SetRedisString: " + err.Error())
	}
	return err
}

func Get(key string) (string, error) {
	attemptString, err := Client.Get(key).Result()
	if err != nil {
		log.LogDebug("Error GetRedis: " + err.Error())
	}
	return attemptString, err
}

func SetRedisHash(key string, objectRedis map[string]interface{}, expiration time.Duration) error {
	err := Client.HMSet(key, objectRedis).Err()
	Client.Expire(key, expiration)
	if err != nil {
		log.LogDebug("Error SetRedisHash: " + err.Error())
	}
	return err
}

func Increase(key string, field string) error {
	err := Client.HIncrBy(key, field, 1).Err()
	if err != nil {
		log.LogDebug("Error Increase: " + err.Error())
	}
	return err
}

func Delete(key string) error {
	err := Client.Del(key).Err()
	if err != nil {
		log.LogDebug("Error Delete: " + err.Error())
	}
	return err
}

func GetHash(key string) (map[string]string, error) {
	data, err := Client.HGetAll(key).Result()
	if err != nil {
		log.LogDebug("Error GetHash: " + err.Error())
	}

	return data, err
}

func GetTTLInSecond(key string) (int, error) {
	cd, err := Client.TTL(key).Result()
	inSecond := int(cd.Seconds())

	if err != nil {
		log.LogDebug("Error GetTTLInSecond: " + err.Error())
	}
	return inSecond, err
}

func IncreaseByKey(key string) (int64, error) {
	num, err := Client.Incr(key).Result()
	if err != nil {
		log.LogDebug("Error IncreaseByKey: " + err.Error())
	}

	return num, err
}
