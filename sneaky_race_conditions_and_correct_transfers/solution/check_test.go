package main

import (
	"sync"
	"testing"
)

func TestTransferPreservesTotalBalance(t *testing.T) {
	accounts := []*Account{
		{ID: 1, balance: 100},
		{ID: 2, balance: 100},
		{ID: 3, balance: 100},
	}

	var wg sync.WaitGroup
	for i := 0; i < 900; i++ {
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
	for _, a := range accounts {
		total += a.Balance()
	}
	if total != 300 {
		t.Fatalf("total balance mismatch: got %d want 300", total)
	}
}
