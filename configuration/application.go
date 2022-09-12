package configuration

import (
	"time"
)

type MysqlDB struct {
	Conn         string
	EnableSqlLog bool
}

type MysqlPool struct {
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime time.Duration
}

type RedisDB struct {
	Host     string
	Password string
	Db       int
}

type RedisPool struct {
	Size    int
	Timeout time.Duration
}

type MongoDb struct {
	Conn           string
	DbName         string
	CollectionName string
}

type LogConf struct {
	FilePath      string
	FileName      string
	ConsoleEnable bool
}

type Application struct {
	Server struct {
		Port           uint
		Name           string
		Environment    string
		MaxHeaderBytes int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
	}
	Logger *LogConf
	Data   struct {
		Source []*MysqlDB
		Pool   *MysqlPool
	}
	Redis struct {
		Source []*RedisDB
		Pool   *RedisPool
	}
	RabbitMq struct {
		Address []string
	}
	Kafka struct {
		Brokers []string
	}
	Mongo struct {
		Source []*MongoDb
	}
	XXLJob struct {
		Addr         string
		Port         string
		AccessToken  string
		ExecutorName string
	}
}
