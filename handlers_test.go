package main

import "net/http"
import "net/http/httptest"
import "net/url"
import "testing"
import "strings"

func makeServerContext() serverContext {
    var sc serverContext
    sc.exit = make (chan int)
    sc.shutdown = make (chan int)
    return sc
}


// Test that if we call the hash handler with no reqeust body,
// we get the expected result code.
func TestHashHandlerReturns400WhenNoFormData(t *testing.T) {
    // The last argument is the request body which we set to nil
    req, err := http.NewRequest("POST", "/hash", nil) 
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()

    context := makeServerContext()

    handler := hashHandler{sc:&context}
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusBadRequest{
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

// Test that if we call the hash handler with a request body
// but no password field, we get the expected result code.
func TestHashHandlerReturns400WhenNoPasswordField(t *testing.T) {
    v := url.Values{}
    v.Set("traffic", "heavy") // None of these are "password"
    v.Set("sound", "loud")
    s := v.Encode()

    req, err := http.NewRequest("POST", "/hash", strings.NewReader(s))

    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    if err != nil {
        t.Fatal(err)
    }
    
    rr := httptest.NewRecorder()
    context := makeServerContext()

    handler := hashHandler{sc:&context}
    handler.ServeHTTP(rr, req)

    // Check status code is what we expect
    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

// Test that if we call the hashHandler with the wrong method
// we get the expected result code
func TestHashHandlerReturns404OnGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/hash", nil) 
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()

    context := makeServerContext()

    handler := hashHandler{sc:&context}
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}
