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
)

func HashAndEncodePassword(pw string) string {
    hashed := sha512.Sum512([]byte(pw))
    return base64.StdEncoding.EncodeToString([]byte(hashed[:]))
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        pw := r.Form["password"][0]
        // TODO - check the dimension of the pw array and handle errors 
        time.Sleep(5 * time.Second)
        fmt.Fprintf(w, HashAndEncodePassword(pw))
        return
    } 
    // Only support the methods explicitly enumerated above, otherwise
    // we return NotFound.
    http.NotFound(w,r)
}

func main() {
    port := flag.Int("port", 8080, "The TCP port to listen on for incoming connections.")
    flag.Parse()
    http.HandleFunc("/hash", hashHandler)
    portString := ":" + strconv.Itoa(*port)
    fmt.Println("Listening on ", portString)
    log.Fatal(http.ListenAndServe(portString, nil))
}
