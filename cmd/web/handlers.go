package main

import (
	"fmt"
	"net/http"
	"strconv"
	"tarala/snippetbox/pkg/forms"
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

	//limit the request body size to 4 MiB
	http.MaxBytesReader(w, r.Body, 4096)
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.Form)

	form.Required("title", "content", "expires")
	form.MaxLength("content", 100)
	form.PermittedValues("expires", "1", "7", "365")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created")
	//after pat mux we don't use query params anymore
	//http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		//we create brand new form, so data is nil
		Form: forms.New(nil),
	})
}

func (app application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//better to use post form
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	//creeate new user
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful, Please log in.")

	http.Redirect(w, r, fmt.Sprintf("/user/login"), http.StatusSeeOther)
}

func (app application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	email := form.Get("email")
	password := form.Get("password")

	id, err := app.users.Authenticate(email, password)
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	// Add a flash message to the session to confirm to the user that they've be
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", 303)
}

func ping(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("OK"))
}
