package cache

import (
    "bj-pfd2/pkg/log"
    "fmt"
    "github.com/go-redis/redis"
    "time"
)

type CfgRedis struct {
    Addr   string
    Passwd string
    Db     int
}

// InitRedisClient 初始化连接
func InitRedisClient(c *CfgRedis) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     c.Addr,
        Password: c.Passwd, // no password set
        DB:       c.Db,     // use default DB
    })

    _, err := rdb.Ping().Result()
    if err != nil {
        panic(fmt.Errorf("init redis client failed: %v", err))
    }
    log.Infof("connected to redis server: %s/%v", c.Addr, c.Db)
    return rdb
}

func redisGet(key string) string {
    if key == "" {
        return ""
    }
    value, err := manager.redisClient.Get(key).Result()
    if err != nil {
        return ""
    } else {
        return value
    }
}

func redisSet(key string, value string, expiration time.Duration) error {
    if value == "" || key == "" {
        return nil
    }
    //log.Infof("Set [%s] in cache", key)
    return manager.redisClient.Set(key, value, expiration).Err()
}
