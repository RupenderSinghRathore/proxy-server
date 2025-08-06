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
	port := flag.String("port", ":4000", "Port to listen on")
	target := flag.String("url", "http://localhost:8080", "Destination Url")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info(*target)
	app := application{
		logger: logger,
		target: *target,
	}
	srv := http.Server{
		Addr:     *port,
		Handler:  app.newRouter(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	app.logger.Info("server started..", "port", *port)
	err := srv.ListenAndServe()
	app.logger.Error(err.Error())
	os.Exit(1)
}
