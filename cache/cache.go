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
	defaultExpiration time.Duration // default cache lifetime
	cleanupInterval   time.Duration // interval at which the cache is cleared
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	var cache = &Cache{
		items:             make(map[string]Item),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.RunCleanExpiredItems()
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


func (o *Cache) Del(key string) error {
	o.Lock()
	defer o.Unlock()

	_, ok := o.items[key]
	if !ok {
		return errors.New("Key not fount")
	}

	delete(o.items, key)

	return nil
}

func (o *Cache) Count() int {
	o.RLock()
	var c = len(o.items)
	o.RUnlock()
	return c
}

// Rename - renamed key
func (o *Cache) Rename(oldKey, newKey string) error {
	o.Lock()
	defer o.Unlock()

	item, ok := o.items[oldKey]
	if !ok {
		return errors.New("Key not found")
	}

	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		return errors.New("This item - expired")
	}

	o.items[newKey] = item

	return nil
}

func (o *Cache) RunCleanExpiredItems() {
	go o.cleanExpiredItems()
}

func (o *Cache) cleanExpiredItems() {
	var t = time.NewTimer(o.cleanupInterval)
	defer t.Stop()

	for {
		<-t.C

		if o.items == nil {
			return
		}

		var keys = o.getExpiredKeys()
		if len(keys) == 0 {
			continue
		}

		o.clearItems(keys)
	}
}

// getExpiredKeys returned all expired keys
func (o *Cache) getExpiredKeys() []string {
	var keys []string

	o.RLock()
	for key, item := range o.items {
		if time.Now().Unix() > item.Expiration && item.Expiration > 0 {
			keys = append(keys, key)
		}
	}
	o.Unlock()

	return keys
}

// clearItems clear all expired items from cache
func (o *Cache) clearItems(keys []string) {
	o.Lock()
	for k := range o.items {
		delete(o.items, k)
	}
	o.Unlock()
}
