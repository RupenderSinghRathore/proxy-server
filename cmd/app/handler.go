package main

import (
	"fmt"
	"net/http"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("path: %s\n", r.URL.Path)
}
