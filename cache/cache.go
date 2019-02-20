package cache

import (
	"errors"
	"sync"
	"time"
)

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64 // lifetime
}

type Cache struct {
	sync.RWMutex
	items             map[string]Item
	defaultExpiration time.Duration // default cache item lifetime
	cleanupInterval   time.Duration // interval at which the cache is cleared
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	var cache = &Cache{
		items:             make(map[string]Item),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		go cache.runCleanExpiredItems()
	}

	return cache
}

func (o *Cache) Set(key string, value interface{}, durationExpired time.Duration) {
	var expiration int64

	if durationExpired <= 0 {
		durationExpired = o.defaultExpiration
	}

	if durationExpired > 0 {
		expiration = time.Now().Add(durationExpired).Unix()
	}

	o.Lock()
	o.items[key] = Item{
		Expiration: expiration,
		Created:    time.Now(),
		Value:      value,
	}
	o.Unlock()
}

func (o *Cache) Get(key string) (interface{}, bool) {
	o.RLock()
	item, ok := o.items[key]
	o.RUnlock()
	if !ok {
		return nil, false
	}

	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		return nil, false
	}

	return item, true
}

func (o *Cache) Del(key string) {
	o.Lock()
	delete(o.items, key)
	o.Unlock()
}

// Exist - check exist item with specified key
func(o *Cache) Exist(key string) bool {
	o.RLock()
	_, ok := o.items[key]
	o.RUnlock()
	return ok
}

func (o *Cache) Count() int {
	o.RLock()
	var c = len(o.items)
	o.RUnlock()
	return c
}

// Rename - renamed key
func (o *Cache) Rename(oldKey, newKey string) error {
	o.RLock()
	item, ok := o.items[oldKey]
	o.Unlock()
	if !ok {
		return errors.New("Key not found")
	}

	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		return errors.New("This item - expired")
	}
	o.Lock()
	o.items[newKey] = item
	o.Unlock()

	return nil
}

// IsCacheExpired - check cache on expired
func (o *Cache) IsCacheExpired() bool {
	o.Lock()
	for _, item := range o.items {
		if item.Expiration == 0 {
			return false
		}

		if time.Now().Unix() < item.Expiration && item.Expiration > 0 {
			return false
		}
	}
	o.Unlock()

	return true
}

func (o *Cache) FlushAll() {
	o.Lock()
	o.items = map[string]Item{}
	o.Unlock()
	return
}

func (o *Cache) runCleanExpiredItems() {
	var t = time.NewTimer(o.cleanupInterval)
	defer t.Stop()

	for {
		<-t.C

		if o.items == nil {
			return
		}

		o.clearItems()
	}
}

// clearItems clear all expired items from cache
func (o *Cache) clearItems() {
	o.Lock()
	for key, item := range o.items {
		if time.Now().Unix() > item.Expiration && item.Expiration > 0 {
			delete(o.items, key)
		}
	}
	o.Unlock()
}
