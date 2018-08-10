package main 

import (
    "fmt"
    "log"
    "crypto/sha512"
    "encoding/base64"
    "encoding/json"
    "net/http"
    "time"
    "strconv"
    "flag"
    "sync"
    "context"
)

var srv http.Server

// A utility function to take a sha512 hash of a password
// and base64 encode the result.
//
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

// This function updates the average with the provided
// value and also increments the counter.
func (avg *safeAverager) updateAverage(usecs int64) {
    avg.mux.Lock()
    var totalTime float64 = avg.avgUsecs * float64(avg.count) + float64(usecs)
    avg.count += 1
    avg.avgUsecs = totalTime / float64(avg.count)
    avg.mux.Unlock()
}

func (avg *safeAverager) getValues() (count int64, avgUsecs float64) {
    avg.mux.Lock()
    defer avg.mux.Unlock()
    return avg.count, avg.avgUsecs
}

type serverContext struct {
    exit chan int          // Channel to signal main that work is done, time to exit
    shutdown chan int      // Channel to signal that a shutdown request was received
    averager safeAverager  // Used to track stats in a thread safe manner.
}

// A function to initialize a serverContext
func (sc *serverContext) init() {
    sc.exit = make(chan int)
    sc.shutdown = make(chan int)
}

/////////////////////////////////////////
// Define  a type for each HTTP handler / 
/////////////////////////////////////////
type hashHandler struct {
    sc *serverContext
}

func (handler *hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        start := time.Now()
        if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if pw,ok := r.Form["password"]; ok && len(pw) == 1 {
            time.Sleep(5 * time.Second)
            fmt.Fprintf(w, HashAndEncodePassword(pw[0]))
            elapsed := time.Since(start)
            // Track stats:
            handler.sc.averager.updateAverage(elapsed.Nanoseconds()/1000) // Dividing by 1000 to convert ns to us
            return
        } else {
            // No password in request, or more than 1. 
            http.Error(w, "Invalid Form Data", http.StatusBadRequest)
            return
        }
    } else {
        // Only support POST
        // otherwise we respond with NotFound.
        http.NotFound(w,r)
    }
}

type shutdownHandler struct {
    sc *serverContext
}

func (handler *shutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Shutdown request acknowledged.")

    // Write to the shutdown channel which will cause another 
    // goroutine to call srv.Shutdown()
    handler.sc.shutdown <- 1

    return
}

type statsHandler struct {
    sc *serverContext
}

func (handler *statsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Grab the stats
        total, avg := handler.sc.averager.getValues()
        averageAsInt := int64(avg) // The averager returns float64, truncate it to int64 for output 

        // Define a structure for our JSON output
        type Stats struct {
            Total int64 `json:"total"`
            Avg   int64 `json:"average"`
        }
        o := Stats{Total:total, Avg:averageAsInt}

        // Create JSON object from the structure
        if output, err := json.Marshal(o); err == nil {
            w.Header().Set("Content-Type", "application/json")
            fmt.Fprintf(w, "%s", output)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        // Only support GET 
        // otherwise we respond with NotFound.
        http.NotFound(w,r)
    }
}
/////////////////////////////////////////
//      End of HTTP handler section     /
/////////////////////////////////////////

func doShutdownWhenChannelSignaled(shutdown chan int, exit chan int) {
    _ = <- shutdown // Block until a shutdown request is received    
    srv.Shutdown(context.Background()) // This blocks until any pending HTTP handlers complete 
    exit <- 1 // Signal main that it's ok to exit now
}

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
