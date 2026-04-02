package main

import (
	"fmt"
	"sync"
	"time"
)

type call struct {
	wg  sync.WaitGroup
	val string
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (string, error)) (string, error, bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}

	c := &call{}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err, false
}

func main() {
	var (
		g       Group
		wg      sync.WaitGroup
		callsMu sync.Mutex
		calls   int
	)

	expensive := func() (string, error) {
		callsMu.Lock()
		calls++
		callsMu.Unlock()

		time.Sleep(100 * time.Millisecond)
		return "value-for-user:42", nil
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			val, err, shared := g.Do("user:42", expensive)
			fmt.Printf("worker=%d val=%q err=%v shared=%v\n", id, val, err, shared)
		}(i)
	}

	wg.Wait()
	fmt.Println("expensive function calls:", calls)
}
