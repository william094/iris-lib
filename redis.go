package iris_lib

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

func InitRedis(cfg *RedisDB, pool *RedisPool) *redis.Client {
	return createRedisConnection(cfg.Host, cfg.Password, cfg.Db, pool.Size, pool.Timeout)
}

func createRedisConnection(addr string, pwd string, db, size int, timeout time.Duration) *redis.Client {
	redisConnection := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    pwd,
		DB:          db,
		PoolSize:    size,
		PoolTimeout: time.Second * timeout,
	})
	if _, err := redisConnection.Ping().Result(); err != nil {
		SystemLogger.Error("redis init failed", zap.String("redis addr", addr), zap.Error(err))
		panic(err)
	} else {
		SystemLogger.Info("redis init success", zap.String("redis addr", addr))
		return redisConnection
	}
}
