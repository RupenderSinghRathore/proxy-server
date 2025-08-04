package main

import "net/http"

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", redirect)

	return mux
}
