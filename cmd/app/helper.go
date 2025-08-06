package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

func makeRequest(method, url string, r *http.Request) (*http.Response, error) {
	req, err := http.NewRequest(method, url, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	res.Header.Add("X-Cache", "MISS")
	return res, err
}

func (app *application) assertServerError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.logger.Error(err.Error())
		fmt.Fprintf(os.Stdout, "Stack traceback: %s", debug.Stack())
	}
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			app.assertServerError(w, fmt.Errorf("%s", err))
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
