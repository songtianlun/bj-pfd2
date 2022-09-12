package cache

import (
	"bj-pfd2/com/log"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
	"time"
)

type cManager struct {
	enable      bool
	t           string
	addr        string
	passwd      string
	db          int
	redisClient *redis.Client
	memCache    *cache.Cache
}

var manager *cManager

// Init 初始化 Cache
//   enable: 是否启用
//   t: 类型，redis/memory
//   addr: redis 地址
//   passwd: redis 密码
//   db: redis db
func Init(enable bool, t string, addr string, passwd string, db int) {
	manager = &cManager{
		enable: enable,
		t:      t,
		addr:   addr,
		passwd: passwd,
		db:     db,
	}
	if enable {
		switch t {
		case "redis":
			manager.redisClient = InitRedisClient(&CfgRedis{
				Addr:   addr,
				Passwd: passwd,
				Db:     db,
			})
		case "memory":
			manager.memCache = initMemoryCache()
		default:
			panic(fmt.Errorf("unknown cache type: %s", t))
		}
	}
}

func Get(key string) string {
	if key == "" || manager == nil || !manager.enable {
		return ""
	}
	switch manager.t {
	case "redis":
		log.DebugF("Get cache [%s] from redis", key)
		return redisGet(key)
	case "memory":
		log.DebugF("Get cache [%s] from memory", key)
		return memoryGet(key)
	default:
		return ""
	}

}

func Set(key string, value string, expiration time.Duration) error {
	if value == "" || key == "" || manager == nil || !manager.enable {
		return nil
	}
	switch manager.t {
	case "redis":
		log.DebugF("Set cache [%s] to redis", key)
		return redisSet(key, value, expiration)
	case "memory":
		log.DebugF("Set cache [%s] to memory", key)
		return memorySet(key, value, expiration)
	default:
		return fmt.Errorf("unknown cache type: %s", manager.t)
	}
}
