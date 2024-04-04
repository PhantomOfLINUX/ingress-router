package main

import (
	"log"
	"net/http"

	"github.com/PhantomOfLINUX/ingress-router/internal/proxy"
)

func main() {
    http.HandleFunc("/", proxy.HandleProxy)

    log.Println("Proxy server is running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}