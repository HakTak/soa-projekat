package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "{\"status\":\"gateway ok\"}")
    })

    fmt.Println("Gateway running on port 8080...")
    http.ListenAndServe(":8080", nil)
}

