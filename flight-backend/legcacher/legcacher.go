package legcacher

import (
	"log"
	"sync"
	"time"

	"github.com/frrad/flight-search/flight-backend/legfinder"
)

type LegCacher struct {
	finder legfinder.LegFinder
	cache  legCache
	ttl    time.Duration
}

type legCache struct {
	sync.RWMutex
	data map[string]cacheValue
}

type cacheValue struct {
	timestamp time.Time
	legs      []legfinder.Leg
}

func NewLegCacher(findr legfinder.LegFinder, ttl time.Duration) *LegCacher {
	newCache := legCache{data: map[string]cacheValue{}}

	return &LegCacher{
		finder: findr,
		cache:  newCache,
		ttl:    ttl,
	}
}

func (c *legCache) set(key string, value []legfinder.Leg) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = cacheValue{
		legs:      value,
		timestamp: time.Now(),
	}
}

func (c *legCache) get(key string) (cacheValue, bool) {
	c.RLock()
	defer c.RUnlock()
	ans, ok := c.data[key]
	return ans, ok
}

func (lc *LegCacher) fillCache(query legfinder.LegSpec) error {
	result, err := lc.finder.Find(query)
	if err != nil {
		return err
	}

	lc.cache.set(query.Hash(), result)
	return nil
}

func (lc *LegCacher) Find(query legfinder.LegSpec) ([]legfinder.Leg, error) {
	cacheVal, ok := lc.cache.get(query.Hash())
	if ok && time.Now().Sub(cacheVal.timestamp) < lc.ttl {
		log.Println("Cache hit.")
		return cacheVal.legs, nil
	}

	log.Println("Cache miss!")
	err := lc.fillCache(query)
	if err != nil {
		return nil, err
	}
	return lc.Find(query)
}
