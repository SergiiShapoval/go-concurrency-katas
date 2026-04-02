package main

import (
	"fmt"
	"sync"
)

type Service struct {
	once  sync.Once
	cond  *sync.Cond
	ready bool
}

func NewService() *Service {
	return &Service{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Service) Start() {
	// TODO: initialize once and broadcast readiness
}

func (s *Service) WaitUntilReady() {
	// TODO: wait on the condition variable
}

func main() {
	service := NewService()

	var waiters sync.WaitGroup
	for i := 1; i <= 4; i++ {
		id := i
		waiters.Add(1)
		go func() {
			defer waiters.Done()
			service.WaitUntilReady()
			fmt.Println("waiter released:", id)
		}()
	}

	var starters sync.WaitGroup
	for i := 1; i <= 3; i++ {
		id := i
		starters.Add(1)
		go func() {
			defer starters.Done()
			service.Start()
			fmt.Println("start called by goroutine:", id)
		}()
	}

	starters.Wait()
	waiters.Wait()
}
