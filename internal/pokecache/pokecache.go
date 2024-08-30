package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cachedData map[string]CacheEntry
	mutex      sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{cachedData: make(map[string]CacheEntry)}
	go newCache.reapLoop(interval)
	return &newCache
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for _ = range ticker.C {
		for key, cacheEntry := range cache.cachedData {
			if time.Since(cacheEntry.createdAt) > interval {
				cache.mutex.Lock()
				delete(cache.cachedData, key)
				cache.mutex.Unlock()
			}
		}
	}
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.cachedData[key] = CacheEntry{createdAt: time.Now(), val: val}
	fmt.Printf("Added %v to cache\n", key)
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cacheEntry, ok := cache.cachedData[key]
	if !ok {
		return nil, false
	}
	fmt.Printf("Got %v from cache\n", key)
	return cacheEntry.val, true
}
