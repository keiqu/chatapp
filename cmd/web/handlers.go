package main

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/justinas/nosurf"

	"github.com/lazy-void/chatapp/internal/models"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

	errors := make(map[string]string)
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	email := r.PostForm.Get("email")

	if strings.TrimSpace(username) == "" {
		errors["username"] = "This field cannot be empty."
	}

	if strings.TrimSpace(email) == "" {
		errors["email"] = "This field cannot be empty."
	} else if !EmailRX.MatchString(email) {
		errors["email"] = "Incorrect email address."
	}

	if strings.TrimSpace(password) == "" {
		errors["password"] = "This field cannot be empty."
	} else if len(password) < 8 {
		errors["password"] = "Password is too short."
	} else if len(password) > 20 {
		errors["password"] = "Password is too long."
	}

	if len(errors) != 0 {
		app.render(w, r, "signup.page.gohtml", templateData{
			Errors: errors,
			Form:   r.PostForm,
		})
		return
	}

	err = app.users.Insert(username, email, password)
	if err == models.ErrDuplicateEmail {
		errors["email"] = "Email is already in use."
		app.render(w, r, "signup.page.gohtml", templateData{
			Errors: errors,
			Form:   r.PostForm,
		})
		return
	} else if err == models.ErrDuplicateUsername {
		errors["username"] = "Username is already taken."
		app.render(w, r, "signup.page.gohtml", templateData{
			Errors: errors,
			Form:   r.PostForm,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	s, _ := app.sessions.Get(r, userSessionKey)
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

	errors := make(map[string]string)
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	if strings.TrimSpace(email) == "" {
		errors["email"] = "This field cannot be empty."
	}
	if strings.TrimSpace(password) == "" {
		errors["password"] = "This field cannot be empty."
	}

	if len(errors) != 0 {
		app.render(w, r, "login.page.gohtml", templateData{
			Errors:    errors,
			Form:      r.PostForm,
			CSRFToken: nosurf.Token(r),
		})
		return
	}

	s, err := app.sessions.Get(r, userSessionKey)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, err := app.users.Authenticate(email, password)
	if err == models.ErrNoRecord {
		s.AddFlash("Account with such email doesn't exist.", "error_flash")

		app.render(w, r, "login.page.gohtml", templateData{
			Errors: errors,
			Form:   r.PostForm,
		})
		return
	} else if err == models.ErrInvalidPassword {
		s.AddFlash("Incorrect password.", "error_flash")

		app.render(w, r, "login.page.gohtml", templateData{
			Errors: errors,
			Form:   r.PostForm,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	s.Values["userID"] = userID
	err = s.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	s, err := app.sessions.Get(r, userSessionKey)
	if err != nil {
		app.serverError(w, err)
		return
	}

	delete(s.Values, "userID")
	err = s.Save(r, w)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
