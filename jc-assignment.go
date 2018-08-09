package main 

import (
    "fmt"
    "log"
    "crypto/sha512"
    "encoding/base64"
    "net/http"
    "time"
    "strconv"
    "flag"
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

type serverContext struct {
    exit chan int // Channel to signal main that work is done, time to exit
    shutdown chan int  // Channel to signal that a shutdown request was received
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
//  
// Define  a type for each HTTP handler 
//
//
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
type hashHandler struct {
    sc *serverContext
}

func (handler *hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if pw,ok := r.Form["password"]; ok && len(pw) == 1 {
            time.Sleep(5 * time.Second)
            fmt.Fprintf(w, HashAndEncodePassword(pw[0]))
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

    log.Println("Shutdown handler OUT")
     
    return
    
}

func doShutdownWhenChannelSignaled(shutdown chan int, exit chan int) {
    _ = <- shutdown // Block until a shutdown request is received    
    log.Println("Calling srv.Shutdown()... ###")
    srv.Shutdown(context.Background()) // This blocks until any pending HTTP handlers complete 
    log.Println("Call to srv.Shutdown() returned *****")
    exit <- 1 // Signal main that it's ok to exit now
}

func main() {
    port := flag.Int("port", 8080, "The TCP port to listen on for incoming HTTP connections.")
    flag.Parse()

    // Create a serverContext to be shared by the http handler goroutines
    var theContext serverContext

    // Initialize the channels
    theContext.exit = make(chan int)
    theContext.shutdown = make(chan int)
   
    // Launch goroutine to wait for shutdown to be signaled
    go doShutdownWhenChannelSignaled(theContext.shutdown, theContext.exit) 

    // Register the handlers. They all wrap the same serverContext.
    http.Handle("/hash", &hashHandler{sc:&theContext})
    http.Handle("/shutdown", &shutdownHandler{sc:&theContext})

    // Setup server to listen on specified port
    portString := ":" + strconv.Itoa(*port)
    log.Println("Listening on ", portString)
    srv.Addr = portString 

    // Start listening for incoming connections
    // This blocks until srv.Shutdown() is called in doShutdownWhenChannelSignaled
    retval := srv.ListenAndServe()
    log.Println("ListenAndServe returned ", retval)

    log.Println("main(): Waiting for shutdown handler to signal done...")
    // Wait for the shutdown handler to signal it's done.
    _ = <- theContext.exit // Blocks until signaled by doShutdownWhenChannelSignaled 
    log.Println("Main OUT")
}
