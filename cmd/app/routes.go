package main

import "net/http"

func (app *application) newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.redirect)

	return mux
}
