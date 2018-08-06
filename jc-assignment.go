package main 


import (
    "fmt"
    "crypto/sha512"
    "encoding/base64"
)

func HashAndEncodePassword(pw string) string {
    hashed := sha512.Sum512([]byte(pw))
    return base64.StdEncoding.EncodeToString([]byte(hashed[:]))
}

func main() {
    password := "angryMonkey"
    output := HashAndEncodePassword(password)
    fmt.Println(output)
}
