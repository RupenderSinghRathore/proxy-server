package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RupenderSinghRathore/proxy-server/internal/response"
)

func (app *application) redirect(w http.ResponseWriter, r *http.Request) {
	resMap := make(map[string]response.Resp)
	ok := app.cachedResponse(w, r, resMap)
	if ok {
		return
	}
	url := fmt.Sprint(app.target + r.URL.Path)
	res, err := app.makeRequest(url, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}
	resMap[r.URL.Path] = response.Resp{
		Header: res.Header,
		Body:   body,
	}
	data, err := json.Marshal(resMap)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "data", data)
	for key, vals := range res.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.Header().Set("X-Cache", "MISS")
	w.Write(body)
}
