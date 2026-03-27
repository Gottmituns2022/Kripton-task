package safemap

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

const year = 1879 // Albert Einstein birth year

func TestSafeMap(t *testing.T) {
	sm := New()

	var wg sync.WaitGroup
	wg.Add(4)

	for workerID := 0; workerID < 4; workerID++ {
		go func(id int) {
			defer wg.Done()

			keys := make([]int, 0, year)
			for key := 1; key <= year; key++ {
				if key%4 == id {
					continue
				}
				keys = append(keys, key)
			}

			r := rand.New(rand.NewSource(100 + int64(id)))
			r.Shuffle(len(keys), func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})

			for _, key := range keys {
				value := sm.GetOrCreate(key)
				atomic.AddInt64(value, 1)
			}
		}(workerID)
	}

	wg.Wait()

	for key := 1; key <= year; key++ {
		got, ok := sm.Load(key)
		if !ok {
			t.Fatalf("key %d is missing", key)
		}
		if got != 3 {
			t.Fatalf("key %d: got %d, want 3", key, got)
		}
	}

	accesses, additions := sm.Stats()
	if accesses != year*3 {
		t.Fatalf("accesses: got %d, want %d", accesses, year*3)
	}
	if additions != year {
		t.Fatalf("additions: got %d, want %d", additions, year)
	}
}
