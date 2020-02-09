package util

import (
	"sync"
	"time"
)

type CacheItem struct {
	Content     interface{}
	ExpiredTime int64
}

var (
	cache      = map[string]*CacheItem{}
	cacheMutex sync.Mutex
)

func GetFromCache(id string) interface{} {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	item, ok := cache[id]
	if !ok {
		return nil
	}

	if item.ExpiredTime < time.Now().Unix() {
		delete(cache, id)
		return nil
	}

	return item.Content
}

func SaveToCache(id string, content interface{}, expiredTime int64) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	item := &CacheItem{
		Content:     content,
		ExpiredTime: expiredTime,
	}
	cache[id] = item
}

func CleanCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for id, _ := range cache {
		delete(cache, id)
	}
}

func CleanExpiredCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for id, item := range cache {
		if item.ExpiredTime < time.Now().Unix() {
			delete(cache, id)
		}
	}
}
