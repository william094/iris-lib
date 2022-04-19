package iris_lib

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Init(db *MongoDb) *mongo.Client {
	op := options.Client().ApplyURI(db.Conn).
		SetConnectTimeout(10 * time.Second).
		SetHeartbeatInterval(time.Minute)
	if client, err := mongo.Connect(context.TODO(), op); err != nil {
		panic(err)
	} else if errs := client.Ping(context.TODO(), nil); err != nil {
		panic(errs)
	} else {
		SystemLogger.Info("mongodb connect success")
		return client
	}
}
