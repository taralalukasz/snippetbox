package main

import (
	"fmt"
	"net/http"
	"strconv"
	"tarala/snippetbox/pkg/models"
)

func (app application) home(w http.ResponseWriter, r *http.Request) {
	//we can now remove it, because pat mux has this feature ootb
	//if r.URL.Path != "/" {
	//	app.notFound(w)
	//	return
	//}

	allSnippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: allSnippets}
	app.render(w, r, "home.page.tmpl", data)
}

func (app application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	//here we pass information from backend to frontend
	data := &templateData{Snippet: snippet}
	app.render(w, r, "show.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//after using pat mux  this is redundant / superflous
	//if r.Method != "POST" {
	//	w.Header().Set("Allow", "POST")
	//	app.clientError(w, http.StatusMethodNotAllowed)
	//	return
	//}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}
	//after pat mux we don't use query params anymore
	//http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
