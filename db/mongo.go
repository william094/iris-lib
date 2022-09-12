package db

import (
	"context"
	"github.com/william094/iris-lib/configuration"
	"github.com/william094/iris-lib/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InitMongo(db *configuration.MongoDb) *mongo.Client {
	op := options.Client().ApplyURI(db.Conn).
		SetConnectTimeout(10 * time.Second).
		SetHeartbeatInterval(time.Minute)
	if client, err := mongo.Connect(context.TODO(), op); err != nil {
		panic(err)
	} else if errs := client.Ping(context.TODO(), nil); err != nil {
		panic(errs)
	} else {
		logx.SystemLogger.Info("mongodb connect success")
		return client
	}
}
