package main

import (
	"errors"
	"net/http"

	"github.com/lazy-void/chatapp/internal/forms"

	"github.com/justinas/nosurf"

	"github.com/lazy-void/chatapp/internal/models"
)

type templateData struct {
	Username     string
	CSRFToken    string
	Form         forms.Form
	SuccessFlash string
	ErrorFlash   string
	Errors       map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "chat.page.gohtml", templateData{})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.gohtml", templateData{})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("username")
	form.ContainsOnlyAllowedChars("username")
	form.MaxLength("username", 50)

	form.Required("email")
	form.MatchesPattern("email", forms.EmailRX)

	form.Required("password")
	form.MinLength("password", 8)
	form.MaxLength("password", 20)

	if !form.Valid() {
		app.render(w, r, "signup.page.gohtml", templateData{
			Form: form,
		})
		return
	}

	err = app.users.Insert(form.Get("username"), form.Get("email"), form.Get("password"))
	switch {
	case errors.Is(err, models.ErrDuplicateEmail):
		form.Errors["email"] = "Email is already in use."
		app.render(w, r, "signup.page.gohtml", templateData{
			Form: form,
		})
		return
	case errors.Is(err, models.ErrDuplicateUsername):
		form.Errors["username"] = "Username is already taken."
		app.render(w, r, "signup.page.gohtml", templateData{
			Form: form,
		})
		return
	case err != nil:
		app.serverError(w, err)
		return
	}

	s, _ := app.sessions.Get(r, sessionKey)
	s.AddFlash("Your signup was successful. Please log in.", "success_flash")
	err = s.Save(r, w)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", templateData{CSRFToken: nosurf.Token(r)})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("email")
	form.MatchesPattern("email", forms.EmailRX)

	form.Required("password")

	if !form.Valid() {
		app.render(w, r, "login.page.gohtml", templateData{
			Form:      form,
			CSRFToken: nosurf.Token(r),
		})
		return
	}

	s, err := app.sessions.Get(r, sessionKey)
	if err != nil {
		app.serverError(w, err)
		return
	}

	username, err := app.users.Authenticate(form.Get("username"), form.Get("password"))
	switch {
	case errors.Is(err, models.ErrNoRecord):
		s.AddFlash("Account with such email doesn't exist.", "error_flash")

		app.render(w, r, "login.page.gohtml", templateData{
			Form: form,
		})
		return
	case errors.Is(err, models.ErrInvalidPassword):
		s.AddFlash("Incorrect password.", "error_flash")

		app.render(w, r, "login.page.gohtml", templateData{
			Form: form,
		})
		return
	case err != nil:
		app.serverError(w, err)
		return
	}

	s.Values[usernameKey] = username
	err = s.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.deleteCookieAuthentication(w, r)
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
