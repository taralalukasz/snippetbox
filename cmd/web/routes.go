package main

import (
	"github.com/bmizerany/pat" // New import
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	//APPROACH 1 - USED STANDARD MUX SERVER WHICH DOESN'T SUPPORT MORE HTTP METHODS ON THE SAME ENDPOINT
	//standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet", app.showSnippet)
	//mux.HandleFunc("/snippet/create", app.createSnippet)
	//
	//fileServer := http.FileServer(http.Dir(".ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	//
	////MIDDLEWARE functions
	////add headers to mux handler, before any request is sent to mux
	////return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	//
	////alternative to above -use alice to simplify syntax
	//return standardMiddleware.Then(mux)

	//APPROACH 2 - USE PAT MUX SO WE CAN EASILY MAKE TWO ENDPOINTS snippet/create
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//we use session Handler middleware from golangcollege lib
	//it has to wrap  each route we want to enable session for
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	//you have to register this before :id, because  mux tries to guess the endpoint. We need to go from general to specific
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
