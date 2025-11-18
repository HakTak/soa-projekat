package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "{\"status\":\"AUTH is ok\"}")
    })

    fmt.Println("Auth running on port 8082...")
    http.ListenAndServe(":8082", nil)
}
