package main

import (
	"fmt"
	"sync"
)

type Config struct {
	Value string
}

var (
	badDone bool
	badCfg  *Config

	goodOnce sync.Once
	goodCfg  *Config
	goodInit = func() *Config {
		return &Config{Value: "hello"}
	}
)

// BadGet shows the classic incorrect idea:
// seeing badDone == true does not guarantee seeing the write to badCfg.
func BadGet() *Config {
	if !badDone {
		badCfg = &Config{Value: "hello"}
		badDone = true
	}
	return badCfg
}

func GoodGet() *Config {
	goodOnce.Do(func() {
		goodCfg = goodInit()
	})
	return goodCfg
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(GoodGet().Value)
		}()
	}
	wg.Wait()
}
