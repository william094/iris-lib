package iris_lib

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func RedisKeyFormat(key string, params ...interface{}) string {
	return fmt.Sprintf(key, params...)
}

//GetValue 获取值
func GetValue(db *redis.Client, key string, data interface{}) error {
	value, err := db.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return err
	}
	if len(value) == 0 {
		return nil
	}
	if err = json.Unmarshal(value, data); err != nil {
		return err
	}
	return nil
}

//SetValue 存储不过期
func SetValue(db *redis.Client, key string, data interface{}) error {
	value, _ := json.Marshal(data)
	if _, err := db.Set(key, value, -1).Result(); err != nil {
		return err
	}
	return nil
}

//SetValueByTime 存储并设置过期时间
func SetValueByTime(db *redis.Client, key string, data interface{}, duration time.Duration) error {
	value, _ := json.Marshal(data)
	if _, err := db.Set(key, value, duration).Result(); err != nil {
		return err
	}
	return nil
}

//HSetValue 存储不过期
func HSetValue(db *redis.Client, key string, filed string, data interface{}) error {
	value, _ := json.Marshal(data)
	if _, err := db.HSet(key, filed, value).Result(); err != nil {
		return err
	}
	return nil
}

func HGetValue(db *redis.Client, key string, filed string, data interface{}) error {
	value, err := db.HGet(key, filed).Bytes()
	if err != nil && err != redis.Nil {
		return err
	}
	if len(value) > 0 {
		if err = json.Unmarshal(value, data); err != nil {
			return err
		}
	}
	return nil
}

func IsMember(db *redis.Client, key string, value interface{}) bool {
	if db.SIsMember(key, value).Val() {
		return true
	}
	return false
}

func DelKey(db *redis.Client, key string) error {
	if _, err := db.Del(key).Result(); err != nil {
		return err
	}
	return nil
}
