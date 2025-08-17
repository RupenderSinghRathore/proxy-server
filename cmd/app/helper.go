package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/RupenderSinghRathore/proxy-server/internal/response"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (app *application) makeRequest(url string, r *http.Request) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	res.Header.Add("X-Cache", "MISS")
	return res, nil
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	app.logger.Error(err.Error())
	fmt.Fprintf(os.Stdout, "Stack traceback: %s", debug.Stack())
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func databaseConn(connStr string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.RequestURI
		)
		app.logger.Info("Request recieved", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}
func (app *application) cachedResponse(w http.ResponseWriter, r *http.Request, resMap map[string]response.Resp) bool {
	data := app.sessionManager.GetBytes(r.Context(), "data")
	json.Unmarshal(data, &resMap)
	path := r.URL.Path
	resp, ok := resMap[path]
	if !ok {
		return false
	}
	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.Header().Set("X-Cache", "HIT")
	w.Write(resp.Body)
	return true
}
