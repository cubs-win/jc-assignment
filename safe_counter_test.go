package main


import "testing"
import "sync"

// A waitgroup to wait for the goroutines to complete before checking results.
var wg sync.WaitGroup

func doSomeIncrements(n int, c *safeCounter) {
    defer wg.Done()

    for i := 0; i < n; i++ {
        c.Increment()
    } 
}

func TestSafeCounter(t *testing.T) {
    var ctr safeCounter
    var numGoroutines int = 25
    var incrementsPerGoroutine int = 100 

    wg.Add(numGoroutines)

    // Launch the goroutines
    for i := 0; i < numGoroutines; i++ {
        go doSomeIncrements(incrementsPerGoroutine, &ctr)
    }

    var expected int = numGoroutines * incrementsPerGoroutine

    wg.Wait() // Wait until the goroutines finish before proceeding.

    total := ctr.Value()
    if int(total) != expected {
        t.Errorf("Counter value %v not equal to expected count %v.", total, expected)
    }
}

