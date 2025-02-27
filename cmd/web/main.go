package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"snippetbox.chapstewie.net/cmd/web/internal/models"
)

type config struct {
	addr      string
	staticDir string
}

var cfg config

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	// config flags
	flag.StringVar(&cfg.addr, "addr", ":4000", "Port to listen to")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	connString := flag.String("connString", "postgres://postgres:password@localhost/snippetbox?sslmode=disable", "PG connection string")
	flag.Parse()

	app := &application{}

	// logger init
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	app.logger = logger

	// db init
	db, err := openDB(*connString)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app.snippets = &models.SnippetModel{DB: db}

	// start server
	app.logger.Info("starting server on ", slog.String("addr", cfg.addr))

	err = http.ListenAndServe(cfg.addr, app.routes())
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
