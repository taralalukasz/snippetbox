package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(".ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//MIDDLEWARE functions
	//add headers to mux handler, before any request is sent to mux
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	//alternative to above -use alice to simplify syntax
	return standardMiddleware.Then(mux)
}
