package main

import (
	"fmt"
	"sync"
	"time"
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
	s.once.Do(func() {
		time.Sleep(100 * time.Millisecond)
		s.cond.L.Lock()
		s.ready = true
		s.cond.Broadcast()
		s.cond.L.Unlock()
		fmt.Println("service initialized")
	})
}

func (s *Service) WaitUntilReady() {
	s.cond.L.Lock()
	for !s.ready {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

func main() {
	service := NewService()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			service.WaitUntilReady()
			fmt.Printf("worker %d started\n", id)
		}(i)
	}

	for i := 0; i < 3; i++ {
		go service.Start()
	}

	wg.Wait()
}
