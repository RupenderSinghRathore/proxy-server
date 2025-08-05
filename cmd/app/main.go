package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	target string
}

func main() {
	addr := flag.String("addr", ":4000", "Port to listen on")
	target := flag.String("url", "http://localhost:8080", "Destination Url")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := application{
		logger: logger,
		target: *target,
	}
	mux := app.newRouter()
	app.logger.Info("server started..", "port", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		app.logger.Error(err.Error())
	}
}
