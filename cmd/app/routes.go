package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) newRouter() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.logRequest)
	mux.Handle("/", dynamic.ThenFunc(app.redirect))

	standard := alice.New(app.recoverPanic, app.setHeaders)
	return standard.Then(mux)
}
