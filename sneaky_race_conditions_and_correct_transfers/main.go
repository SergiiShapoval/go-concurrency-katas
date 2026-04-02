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
	// TODO: return the protected balance
	return 0
}

func Transfer(from, to *Account, amount int) bool {
	// TODO: lock in stable order, move funds, unlock
	return false
}

func main() {
	accounts := []*Account{
		{ID: 1, balance: 100},
		{ID: 2, balance: 100},
		{ID: 3, balance: 100},
	}

	type transfer struct {
		from   int
		to     int
		amount int
	}

	ops := []transfer{
		{from: 0, to: 1, amount: 10},
		{from: 1, to: 2, amount: 15},
		{from: 2, to: 0, amount: 7},
		{from: 0, to: 2, amount: 12},
		{from: 1, to: 0, amount: 5},
		{from: 2, to: 1, amount: 9},
	}

	var wg sync.WaitGroup
	for _, op := range ops {
		op := op
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok := Transfer(accounts[op.from], accounts[op.to], op.amount)
			fmt.Printf("transfer %d -> %d amount=%d ok=%v\n", accounts[op.from].ID, accounts[op.to].ID, op.amount, ok)
		}()
	}
	wg.Wait()

	total := 0
	for _, account := range accounts {
		balance := account.Balance()
		total += balance
		fmt.Printf("account %d balance=%d\n", account.ID, balance)
	}
	fmt.Println("total:", total)
}
