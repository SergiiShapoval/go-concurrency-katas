package main

import (
	"fmt"
	"sync"
)

type Account struct {
	ID      int
	mu      sync.Mutex
	balance int
}

func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func Transfer(from, to *Account, amount int) bool {
	first, second := from, to
	if first.ID > second.ID {
		first, second = second, first
	}

	first.mu.Lock()
	second.mu.Lock()
	defer second.mu.Unlock()
	defer first.mu.Unlock()

	if from.balance < amount {
		return false
	}

	from.balance -= amount
	to.balance += amount
	return true
}

func main() {
	accounts := []*Account{
		{ID: 1, balance: 100},
		{ID: 2, balance: 100},
		{ID: 3, balance: 100},
	}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		from := accounts[i%len(accounts)]
		to := accounts[(i+1)%len(accounts)]
		wg.Add(1)
		go func() {
			defer wg.Done()
			Transfer(from, to, 1)
		}()
	}

	wg.Wait()

	total := 0
	for _, account := range accounts {
		b := account.Balance()
		total += b
		fmt.Printf("account %d => %d\n", account.ID, b)
	}
	fmt.Println("total:", total)
}
