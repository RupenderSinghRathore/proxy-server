package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

type application struct {
	logger         *slog.Logger
	target         string
	baseUrl        string
	sessionManager *scs.SessionManager
}

func main() {
	port := flag.String("port", ":4000", "Port to listen on")
	target := flag.String("url", "http://localhost:8080", "Destination Url")
	connStr := flag.String("conn", "postgres://kami-sama:touka@localhost:5432/proxy_server", "Postgres connection string")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := databaseConn(*connStr)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	app := application{
		logger:         logger,
		target:         *target,
		baseUrl:        *target,
		sessionManager: scs.New(),
	}
	app.sessionManager.Lifetime = 12 * time.Hour
	app.sessionManager.Store = postgresstore.New(db)
	srv := http.Server{
		Addr:     *port,
		Handler:  app.newRouter(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	app.logger.Info("server started..", "port", *port)
	err = srv.ListenAndServe()
	// err := srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	app.logger.Error(err.Error())
	os.Exit(1)
}
