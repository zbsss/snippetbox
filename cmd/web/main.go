package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zbsss/snippetbox/internal/models"
)

type (
	config struct {
		addr      string
		staticDir string
		dsn       string
	}

	application struct {
		logger   *slog.Logger
		cfg      *config
		snippets *models.SnippetModel
	}
)

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := application{
		logger:   logger,
		cfg:      &cfg,
		snippets: &models.SnippetModel{DB: db},
	}

	app.logger.Info("starting server", slog.String("addr", cfg.addr))

	err = http.ListenAndServe(cfg.addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
