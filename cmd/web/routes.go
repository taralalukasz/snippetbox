package main

import "net/http"

func (app application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(".ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//MIDDLEWARE functions
	//add headers to mux handler, before any request is sent to mux
	return secureHeaders(mux)
}