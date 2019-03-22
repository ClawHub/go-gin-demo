package gredis

import (
	"github.com/go-redis/redis"
	"go-gin-demo/pkg/logging"
	"go-gin-demo/pkg/setting"
	"go.uber.org/zap"
)

//redis单机客户端
var client *redis.Client

//redis集群客户端
var redisClusterClient *redis.ClusterClient

//redis配置结构体
type Redis struct {
	Host     string
	Password string
	DB       int
	Cluster  bool
	Hosts    []string
}

var RedisSetting = &Redis{}

//redis设置
func Setup() {

	//配置文件读取
	setting.MapTo("redis", RedisSetting)

	//判断是否为集群配置
	if RedisSetting.Cluster {
		//ClusterClient是一个Redis集群客户机，表示一个由0个或多个底层连接组成的池。它对于多个goroutine的并发使用是安全的。
		redisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Password: RedisSetting.Password,
			Addrs:    RedisSetting.Hosts,
		})
		//Ping
		ping, err := redisClusterClient.Ping().Result()
		logging.AppLogger.Info("Redis cluster Ping", zap.String("ping", ping), zap.Error(err))

	} else {
		//Redis客户端，由零个或多个基础连接组成的池。它对于多个goroutine的并发使用是安全的。
		//更多参数参考Options结构体
		client = redis.NewClient(&redis.Options{
			Addr:     RedisSetting.Host,
			Password: RedisSetting.Password, // no password set
			DB:       RedisSetting.DB,       // use default DB
		})
		//Ping
		ping, err := client.Ping().Result()
		logging.AppLogger.Info("Redis Ping", zap.String("ping", ping), zap.Error(err))
	}

}
