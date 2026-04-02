package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	m sync.Map
}

func (c *Cache) GetOrCompute(key string, fn func() string) string {
	// TODO: use Load / LoadOrStore to populate the cache
	return ""
}

func computeAll(cache *Cache, keys []string, fn func(string) string) []string {
	// TODO: call GetOrCompute concurrently for all keys and preserve result order
	return nil
}

func main() {
	cache := &Cache{}
	var computes atomic.Int32

	keys := []string{"profile", "settings", "profile", "profile", "settings"}
	results := computeAll(cache, keys, func(key string) string {
		n := computes.Add(1)
		time.Sleep(20 * time.Millisecond)
		return fmt.Sprintf("%s-value-%d", key, n)
	})

	fmt.Println("results:", results)
	fmt.Println("compute calls:", computes.Load())
}
