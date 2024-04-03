package main

import (
	"log"
	"net/http"

	"github.com/PhantomOfLINUX/ingress-router/internal/proxy"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		go proxy.HandleProxy(w, r)
	})

	log.Println("Proxy server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
