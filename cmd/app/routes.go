package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) newRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.redirect)

	standard := alice.New(app.recoverPanic, app.setHeaders)
	return standard.Then(mux)
}
