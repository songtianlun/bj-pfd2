package cache

import (
	"bj-pfd2/com/log"
	"github.com/go-redis/redis"
	"time"
)

type CfgRedis struct {
	Addr   string
	Passwd string
	Db     int
}

// 声明一个全局的rdb变量
var rdb *redis.Client

// InitClient 初始化连接
func InitClient(c *CfgRedis) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Passwd, // no password set
		DB:       c.Db,     // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	log.InfoF("connected to redis server: %s/%v", c.Addr, c.Db)
	return nil
}

// InitClient 初始化连接
//func InitClient() (err error) {
//	rdb = redis.NewClient(&redis.Options{
//		Addr:     "192.168.2.5:6379",
//		Password: "wRVP7Dd+I0ZwO2b+/ljEinygkycFMFVCNglSDkvjKtt1n6wAT9yjU9yVAzjse5y4OFhS+ZTK9UHZ00MsDlnTGU1e3f0M4XcvtM6Ro4LlSf0dxwZuh/LDBxhE", // no password set
//		DB:       0,                                                                                                                          // use default DB
//	})
//
//	_, err = rdb.Ping().Result()
//	if err != nil {
//		return err
//	}
//	return nil
//}

func Get(key string) string {
	if key == "" {
		return ""
	}
	value, err := rdb.Get(key).Result()
	if err != nil {
		return ""
	} else {
		return value
	}
}

func Set(key string, value string) error {
	if value == "" || key == "" {
		return nil
	}
	//log.InfoF("Set [%s] in cache", key)
	return rdb.Set(key, value, 30*time.Minute).Err()
}
