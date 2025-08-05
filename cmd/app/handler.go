package main

import (
	"fmt"
	"io"
	"net/http"
)

func (app *application) redirect(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("path: %s\n", r.URL.Path)
	url := fmt.Sprint(app.target + r.URL.Path)
	method := r.Method
	res, err := makeRequest(method, url, r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	defer res.Body.Close()
	for key, vals := range res.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	io.Copy(w, res.Body)
}
