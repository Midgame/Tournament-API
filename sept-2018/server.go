package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func main() {
    http.HandleFunc("/ping", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
