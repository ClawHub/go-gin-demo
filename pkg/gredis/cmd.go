package gredis

import (
	"github.com/go-redis/redis"
	"time"
)

//set
func Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if RedisSetting.Cluster {
		return redisClusterClient.Set(key, value, expiration)
	} else {
		return client.Set(key, value, expiration)
	}
}

//get
func Get(key string) *redis.StringCmd {
	if RedisSetting.Cluster {
		return redisClusterClient.Get(key)
	} else {
		return client.Get(key)
	}
}

//删除
func Delete(keys ...string) *redis.IntCmd {
	if RedisSetting.Cluster {
		return redisClusterClient.Del(keys...)
	} else {
		return client.Del(keys...)
	}
}

//判断是否存在
func Exists(keys ...string) *redis.IntCmd {
	if RedisSetting.Cluster {
		return redisClusterClient.Exists(keys...)
	} else {
		return client.Del(keys...)
	}
}

//订阅
func Subscribe(channels ...string) *redis.PubSub {
	if RedisSetting.Cluster {
		return redisClusterClient.Subscribe(channels...)
	} else {
		return client.Subscribe(channels...)
	}
}

//发布
func Publish(channel string, message interface{}) *redis.IntCmd {
	if RedisSetting.Cluster {
		return redisClusterClient.Publish(channel, message)
	} else {
		return client.Publish(channel, message)
	}
}
