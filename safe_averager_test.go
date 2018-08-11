package main

import "testing"
import "sync"
import "math/rand"

// A waitgroup to wait for the goroutines to complete before checking results.
var wg sync.WaitGroup

var numGoroutines = 25
var updatesPerGoroutine = 100 

// A channel for the goroutines to send their times through 
var myChan chan int64

// A slice into which the times are collected and saved
// so we can check the average calculated by the averager
var times []int64


// doSomeAverging iterates to repeatedly update the running average
// kept by the safeAverager, using a randomly generated time value.
// The time values are also sent back through a channel
// so that the average can be checked against the value reported by
// the averager.
func doSomeAveraging(numIterations int, averager *safeAverager) {
    defer wg.Done()

    for i := 0; i < numIterations; i++ {
        var usecs = rand.Int63n(5000000) + 1  // Random time between 1 usec and 5 sec
        myChan <- usecs
        averager.updateAverage(usecs)
    } 
}

func collectTimes() {
    for time := range(myChan) {
        if time == -1 {
            return
        }
        times = append(times, time)
    }
}

func TestSafeAverager(t *testing.T) {
    var averager safeAverager
    myChan = make(chan int64)

    wg.Add(numGoroutines) 

    // Start the collectTimes goRoutine
    go collectTimes()

    // Launch the goroutines
    for i := 0; i < numGoroutines; i++ {
        go doSomeAveraging(updatesPerGoroutine, &averager)
    }

    wg.Wait() // Wait until the goroutines finish before proceeding.
 
    // Signal the collectTimes() function to exit
    myChan <- int64(-1) 

    count, avg := averager.getValues()

    var expectedCount = int64(numGoroutines * updatesPerGoroutine)

    var sum int64
    for _, val :=  range times {
        sum += val
    }

    var expectedAverage = float64(sum) / float64(expectedCount)

    if count != expectedCount {
        t.Errorf("Counter value %v not equal to expected count %v.", count, expectedCount)
    }

    if int(expectedAverage) != int(avg) {
        t.Errorf("Average value %v not equal to expected average value %v",
                 int(avg), int(expectedAverage))
    }
}

