package main

import "testing"

func TestHash(t *testing.T) {
    pw := "angryMonkey"
    hashed := HashAndEncodePassword(pw)
    expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
    
    if hashed != expected{
        t.Error("Incorrect hash calculated for password ", pw)
    } 
}

