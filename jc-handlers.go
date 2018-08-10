package main

import (
    "fmt"
    "encoding/json"
    "time"
    "net/http"
)
/////////////////////////////////////////
// Define  a type for each HTTP handler / 
/////////////////////////////////////////

// hashHandler is an object to handle calls to the /hash endpoint
type hashHandler struct {
    sc *serverContext
}

// ServeHTTP on the hashHandler object to implement the Handler interface
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
        }
        // No password in request, or more than 1. 
        http.Error(w, "Invalid Form Data", http.StatusBadRequest)
        return
    }
    // Only support POST
    // otherwise we respond with NotFound.
    http.NotFound(w,r)
}

// shutdownHandler is an object to handle calls to the /shutdown endpoint
type shutdownHandler struct {
    sc *serverContext
}

// ServeHTTP on the shutdownHandler object to implement the Handler interface.
func (handler *shutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Shutdown request acknowledged.")

    // Write to the shutdown channel which will cause another 
    // goroutine to call srv.Shutdown()
    handler.sc.shutdown <- 1

    return
}

// statsHandler is an object to handle calls to the /stats endpoint
type statsHandler struct {
    sc *serverContext
}

// ServeHTTP on the statsHandler object to implement the Handler interface.
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
        return
    }
    // Only support GET 
    // otherwise we respond with NotFound.
    http.NotFound(w,r)
}
/////////////////////////////////////////
//      End of HTTP handler section     /
/////////////////////////////////////////


