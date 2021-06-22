package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/justinas/nosurf"

	"github.com/lazy-void/chatapp/internal/models"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, "chat.page.gohtml", templateData{})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, "signup.page.gohtml", templateData{CSRFToken: nosurf.Token(r)})
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
		app.render(w, "signup.page.gohtml", templateData{Errors: errors, Form: r.PostForm})
		return
	}

	err = app.users.Insert(username, email, password)
	if err == models.ErrDuplicateEmail {
		errors["email"] = "Email is already in use."
		app.render(w, "signup.page.gohtml", templateData{Errors: errors, Form: r.PostForm})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, "login.page.gohtml", templateData{Alert: "Your signup was successful. Please log in."})
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, "login.page.gohtml", templateData{CSRFToken: nosurf.Token(r)})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login processing")
}
