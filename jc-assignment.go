package main 

import (
    "log"
    "crypto/sha512"
    "encoding/base64"
    "net/http"
    "strconv"
    "flag"
    "sync"
    "context"
)

var srv http.Server

// HashAndEncodePassword A utility function to take a sha512 hash of a password
// and base64 encode the result.
func HashAndEncodePassword(pw string) string {
    hashed := sha512.Sum512([]byte(pw))
    return base64.StdEncoding.EncodeToString([]byte(hashed[:]))
}

// A counter and  average response time calculator that's safe to use concurrently
type safeAverager struct {
    avgUsecs float64
    count int64 
    mux sync.Mutex
}

// updateAverage updates the running average with the provided
// value and also increments the counter.
func (avg *safeAverager) updateAverage(usecs int64) {
    avg.mux.Lock()
    var totalTime = avg.avgUsecs * float64(avg.count) + float64(usecs)
    avg.count++
    avg.avgUsecs = totalTime / float64(avg.count)
    avg.mux.Unlock()
}

// getValues returns a tuple (count, avgUsecs) containing the current number
// of hash requests that have been serviced and the average response time in microseconds.
func (avg *safeAverager) getValues() (count int64, avgUsecs float64) {
    avg.mux.Lock()
    defer avg.mux.Unlock()
    return avg.count, avg.avgUsecs
}

// serverContext is a structure that holds state for the server process.
type serverContext struct {
    exit chan int          // Channel to signal main that work is done, time to exit
    shutdown chan int      // Channel to signal that a shutdown request was received
    averager safeAverager  // Used to track stats in a thread safe manner.
}

// init is a convenience method to initialize a serverContext.
func (sc *serverContext) init() {
    sc.exit = make(chan int)
    sc.shutdown = make(chan int)
}


// doShutdownWhenChannelSignaled is a function that will run as a goroutine
// to wait for a signal from the shutdownHandler. When the signal comes,
// it will call Shutdown() on the http.Server object.
func doShutdownWhenChannelSignaled(shutdown chan int, exit chan int) {
    _ = <- shutdown // Block until a shutdown request is received    
    srv.Shutdown(context.Background()) // This blocks until any pending HTTP handlers complete 
    exit <- 1 // Signal main that it's ok to exit now
}

// main is the program entry point
func main() {
    port := flag.Int("port", 8080, "The TCP port to listen on for incoming HTTP connections.")
    flag.Parse()

    // Create a serverContext to be shared by the http handler goroutines
    var theContext serverContext

    // Initialize the context
    theContext.init() 
   
    // Launch goroutine to wait for shutdown to be signaled
    go doShutdownWhenChannelSignaled(theContext.shutdown, theContext.exit) 

    // Register the handlers. They all wrap the same serverContext.
    http.Handle("/hash", &hashHandler{sc:&theContext})
    http.Handle("/shutdown", &shutdownHandler{sc:&theContext})
    http.Handle("/stats", &statsHandler{sc:&theContext})

    // Setup server to listen on specified port
    portString := ":" + strconv.Itoa(*port)
    log.Println("Listening on ", portString)
    srv.Addr = portString 

    // Start listening for incoming connections
    // This blocks until srv.Shutdown() is called in doShutdownWhenChannelSignaled
    retval := srv.ListenAndServe()
    log.Println("ListenAndServe returned ", retval)

    log.Println("main(): Waiting for active work to finish...")
    // Wait for the signal that active work is completed before we exit
    _ = <- theContext.exit // Blocks until signaled by doShutdownWhenChannelSignaled 
    log.Println("main() OUT")
}
