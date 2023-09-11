package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type (
	config struct {
		addr      string
		staticDir string
	}

	application struct {
		logger *slog.Logger
	}
)

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	app := application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(cfg.staticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	app.logger.Info("starting server", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, mux)
	app.logger.Error(err.Error())
	os.Exit(1)
}
