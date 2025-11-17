package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Stakeholders OK")
    })

    fmt.Println("Stakeholders service running on port 8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}

