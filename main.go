package main

import (
	"fmt"
	"net/http"
)

func main() {
	addr := "127.0.0.1:3000"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, router())
}

