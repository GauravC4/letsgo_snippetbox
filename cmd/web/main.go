package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

var cfg config

type application struct {
	logger *slog.Logger
}

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "Port to listen to")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	app := &application{
		logger: logger,
	}

	app.logger.Info("starting server on ", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

}
