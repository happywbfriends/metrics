package main

import (
	"fmt"
	metricsv1 "github.com/happywbfriends/metrics/v1"
	"github.com/jellydator/ttlcache/v3"
	"time"
)

func CacheExample() {
	metrics := metricsv1.NewCacheMetrics()
	cache := newCacheInt("Ints", metrics)
	cache.Set("test", 1)
	v, found := cache.Get("test")
	if found {
		fmt.Printf("VAL in cache: %v", v)
	}
}

type CacheInt interface {
	Set(k string, v int)
	Get(k string) (int, bool)
}

func newCacheInt(name string, metrics metricsv1.CacheMetrics) CacheInt {
	var c = cacheInt{
		name: name,
		cache: ttlcache.New(
			ttlcache.WithTTL[string, int](5*time.Minute),
			ttlcache.WithDisableTouchOnHit[string, int](),
		),
		metrics: metrics,
	}

	// Periodical cleanup
	go c.cache.Start()

	return &c
}

type cacheInt struct {
	name    string
	cache   *ttlcache.Cache[string, int]
	metrics metricsv1.CacheMetrics
}

func (c *cacheInt) Set(k string, v int) {
	c.cache.Set(k, v, ttlcache.DefaultTTL)
	c.metrics.SetSize(c.name, 0, c.cache.Len())
}

func (c *cacheInt) Get(k string) (int, bool) {
	item := c.cache.Get(k)
	if item != nil {
		c.metrics.IncNbReadHit(c.name, 0)
		return item.Value(), true
	}

	c.metrics.IncNbReadMiss(c.name, 0)
	return 0, false
}
