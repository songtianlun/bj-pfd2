package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func initMemoryCache() *cache.Cache {
	return cache.New(5*time.Minute, 10*time.Minute)
}

func memoryGet(key string) string {
	v, ok := manager.memCache.Get(key)
	if ok {
		return v.(string)
	} else {
		return ""
	}
}

func memorySet(key string, value string, expiration time.Duration) error {
	manager.memCache.Set(key, value, expiration)
	return nil
}
