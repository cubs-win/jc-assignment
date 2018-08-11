package main

import (
    "fmt"
    "time"
    "net/http"
    "net/url"
)

var hashUrl = "http://localhost:8080/hash"
var shutdownUrl = "http://localhost:8080/shutdown"

var channel chan int

func doHash(pw string) {
    formData := url.Values {
        "password": {pw},
    }
    resp, err := http.PostForm(hashUrl, formData)     
 
    if err != nil {
        fmt.Println("Error in doHash: ", err)
        channel <- -1
    } else {
        channel <- resp.StatusCode
    }
}

func doShutdown() {
    resp, err := http.Get(shutdownUrl)
    if err != nil {
        fmt.Println("Error in doShutdown: ", err)
        channel <- -1
    } else {
        channel <- resp.StatusCode
    }
}

func main() {
   
    // This is a test to make sure we can open multiple
    // concurrent connections to the jc-assignment server,
    // and that issuing a shutdown doesn't prevent responses
    // from being retrieved from previously submitted requests.

    requestCount := 100

    channel = make (chan int, requestCount + 1) // +1 is for the shutdown

  
    // Launch goroutines
    for i := 0; i < requestCount; i++ {
        go doHash("FruityPebbles1")
    }

    // Sleep just long enough to let the requests be issued before we shutdown
    time.Sleep(100 * time.Millisecond) 

    // Issue shutdown request
    go doShutdown()

    failures := 0
    successes := 0

    for i := 0; i < requestCount + 1; i++ {
        rc := <- channel
        if rc != 200 {
            failures++
        } else {
            successes++
        }
    }

    fmt.Printf("Received %v successful response codes and %v failures.\n", successes, failures)


}
