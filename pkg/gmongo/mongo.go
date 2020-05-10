package gmongo

import (
	"context"
	"diffuser/pkg/logging"
	"diffuser/pkg/setting"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var Client *mongo.Client

type Mongo struct {
	Uri     string
	TimeOut time.Duration
}

var MongoSetting = &Mongo{}

func Setup() {
	//读取配置文件
	setting.MapTo("mongo", MongoSetting)
	var err error
	//链接mongo服务
	Client, err = mongo.Connect(GetContext(), options.Client().ApplyURI(MongoSetting.Uri))
	if err != nil {
		logging.AppLogger.Fatal("open mongodb fail")
	}
	//判断服务是否可用
	err = Client.Ping(GetContext(), readpref.Primary())
	if err != nil {
		logging.AppLogger.Fatal("ping mongodb fail")
	}
}

//公共参数
func GetContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), MongoSetting.TimeOut)
	return
}
