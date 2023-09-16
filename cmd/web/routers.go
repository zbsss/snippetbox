package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(app.cfg.staticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.logRequest(secureHeaders(mux))
}
