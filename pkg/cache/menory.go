package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func initMemoryCache() *cache.Cache {
	return cache.New(5*time.Minute, 10*time.Minute)
}

func memoryGet(c *cache.Cache, key string) string {
	v, ok := c.Get(key)
	if ok {
		return v.(string)
	} else {
		return ""
	}
}

func memorySet(c *cache.Cache, key string, value string, expiration time.Duration) error {
	c.Set(key, value, expiration)
	return nil
}
