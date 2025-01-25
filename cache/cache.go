package cache

import (
	"errors"
	"sync"
	"time"
)

type cacheEntry struct {
	URL      string
	ExpireAt time.Time
}

var cache = struct {
	sync.RWMutex
	store map[string]cacheEntry
}{store: make(map[string]cacheEntry)}

func Get(shortURL string) (string, bool) {
	cache.RLock()
	defer cache.RUnlock()
	entry := cache.store[shortURL]
	if entry.URL == "" || time.Now().After(entry.ExpireAt) {
		return "", false
	}
	return entry.URL, true
}
func Add(shortURL string, originalURL string) error {
	ttl := time.Hour * 12
	cache.Lock()
	defer cache.Unlock()
	if shortURL == "" || originalURL == "" {
		return errors.New("shortURL or originalURL cannot be empty")
	}
	cache.store[shortURL] = cacheEntry{
		URL:      originalURL,
		ExpireAt: time.Now().Add(ttl),
	}
	return nil
}
func Cleanup() {
	cache.Lock()
	defer cache.Unlock()
	for key, entry := range cache.store {
		if time.Now().After(entry.ExpireAt) {
			delete(cache.store, key)
		}
	}
}
func StartCacheCleanup() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			Cleanup()
		}
	}()
}
