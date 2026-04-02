package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	m sync.Map
}

func (c *Cache) GetOrCompute(key string, fn func() string) string {
	if v, ok := c.m.Load(key); ok {
		return v.(string)
	}

	value := fn()
	actual, _ := c.m.LoadOrStore(key, value)
	return actual.(string)
}

func computeAll(cache *Cache, keys []string, fn func(string) string) []string {
	out := make([]string, len(keys))
	var wg sync.WaitGroup
	for i, key := range keys {
		wg.Add(1)
		go func(i int, key string) {
			defer wg.Done()
			out[i] = cache.GetOrCompute(key, func() string { return fn(key) })
		}(i, key)
	}
	wg.Wait()
	return out
}

func main() {
	var (
		cache   Cache
		mu      sync.Mutex
		compute int
	)

	expensive := func(key string) func() string {
		return func() string {
			mu.Lock()
			compute++
			mu.Unlock()
			time.Sleep(50 * time.Millisecond)
			return strings.ToUpper(key)
		}
	}

	keys := []string{"go", "go", "channels", "go", "mutex", "channels"}
	results := computeAll(&cache, keys, func(key string) string { return expensive(key)() })
	for i, key := range keys {
		fmt.Printf("%s => %s\n", key, results[i])
	}
	fmt.Println("compute calls:", compute)
}
