package main

import "sync"

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
	// TODO: share one in-flight call per key
	return "", nil, false
}

func main() {
	// TODO: start many goroutines for the same key and prove fn ran once
}
