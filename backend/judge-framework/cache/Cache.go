package cache

import "sync"

type GlobalCache interface {
	Get(key string) string
	Set(key string, value string)
}

type globalCacheImpl struct {
	storage map[string]string
}

func (c globalCacheImpl) Get(key string) string {
	return c.storage[key]
}

func (c globalCacheImpl) Set(key, value string) {
	c.storage[key] = value
}

var cache *globalCacheImpl
var once sync.Once

func GetGlobalCache() GlobalCache {
	once.Do(func() {
		cache = &globalCacheImpl{
			storage: make(map[string]string),
		}
	})
	return *cache
}
