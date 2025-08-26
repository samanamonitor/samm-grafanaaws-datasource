package cache

import (
	"time"

    "github.com/samana-group/sammaws/pkg/samm"

    "github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type CacheState int
const (
    CACHEFULL CacheState = iota
    CACHEPARTIAL
    CACHEEMPTY
)

type Cache struct {
    expires       time.Time
    Objects       interface{}
    NextToken     *string
    state         CacheState
    cacheDuration time.Duration
}

func NewCache(cacheDuration time.Duration) (*Cache) {
    cacheItem := &Cache{
        expires: time.Now().Add(-5 * time.Minute), 
        state: CACHEEMPTY,
        cacheDuration: cacheDuration,
    }
    return cacheItem
}

func (cache *Cache) Flush() {
    cache.Objects = []interface{}{}
    cache.expires = time.Now().Add(-5 * time.Minute)
    cache.NextToken = nil
    cache.state = CACHEEMPTY
}

func (cache *Cache) Update(se samm.SammElement, lastError error) {
    if cache.state == CACHEFULL {
        return
    }
    cache.Objects = se.Elements()
    cache.expires = time.Now().Add(cache.cacheDuration)
    cache.NextToken = se.NextToken()
    if lastError == nil {
        cache.state = CACHEFULL
    } else {
        cache.state = CACHEPARTIAL
    }
    log.DefaultLogger.Info("Cache refreshed.", "expires", cache.expires.String(), 
        "elements", se.Len(), "state", cache.state, "NextToken", cache.NextToken)
}

func (cache *Cache) IsExpired() bool {
    if cache.expires.Sub(time.Now()) <= 0 * time.Second {
        cache.Flush()
        return true
    }
    return false
}

func (cache Cache) State() CacheState {
    return cache.state
}

func (cache *Cache) IsEmpty() bool {
    return cache.state == CACHEEMPTY
}

func (cache *Cache) IsPartial() bool {
    return (! cache.IsExpired()) && cache.state != CACHEFULL
}

func (cache *Cache) IsValid() bool {
    return (! cache.IsExpired()) && (cache.state == CACHEFULL)
}

type CacheMap struct {
    cacheDuration time.Duration
    data map[string]*Cache
}

func NewCacheMap(cacheDuration time.Duration) CacheMap {
    return CacheMap{
        data: make(map[string]*Cache),
        cacheDuration: cacheDuration,
    }
}

func (cm *CacheMap) Get(serviceKey string) (*Cache) {
    cacheItem, ok := cm.data[serviceKey]
    if ! ok {
        cacheItem = NewCache(cm.cacheDuration)
        cm.data[serviceKey] = cacheItem
        log.DefaultLogger.Info("Cache not initialized", "type", serviceKey)
    }

    if cacheItem.IsExpired() {
        log.DefaultLogger.Info("Cache has expired.")

    }
    switch cacheItem.State() {
    case CACHEFULL:
        log.DefaultLogger.Info("Cache Hit.")
    case CACHEEMPTY:
        log.DefaultLogger.Info("Cache Miss.")
    case CACHEPARTIAL:
        log.DefaultLogger.Info("Cache Partial hit. Need to continue.")
    }

    return cacheItem

}
