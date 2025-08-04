package main

import (
	"log"
	"net/http"
)

func main() {
	mux := newRouter()
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
