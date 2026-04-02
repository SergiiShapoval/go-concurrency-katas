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
		return &Config{
			Value: "val",
		}
	}
)

func BadGet() *Config {
	// TODO: intentionally sketch the broken fast-path version
	return nil
}

func GoodGet() *Config {
	// TODO: implement with sync.Once
	return nil
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
