package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/lazy-void/chatapp/chat"
	"github.com/lazy-void/chatapp/models"

	"github.com/justinas/nosurf"
	"github.com/rs/zerolog/log"
)

const (
	sessionKey  = "user-session"
	usernameKey = "username"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	log.Error().Msg(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (app *Application) addDefaultData(w http.ResponseWriter, r *http.Request, td templateData) templateData {
	if user := app.authenticatedUser(r); user != nil {
		td.Username = user.Username
	}

	td.CSRFToken = nosurf.Token(r)

	s, _ := app.Sessions.Get(r, sessionKey)
	successFlashes := s.Flashes("success_flash")
	if len(successFlashes) > 0 {
		td.SuccessFlash = successFlashes[0].(string)
	}

	errorFlashes := s.Flashes("error_flash")
	if len(errorFlashes) > 0 {
		td.ErrorFlash = errorFlashes[0].(string)
	}

	err := s.Save(r, w)
	if err != nil {
		app.serverError(w, err)
	}

	return td
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td templateData) {
	ts, err := template.ParseFS(htmlTemplates, filepath.Join("templates", name), "templates/base.layout.gohtml")
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf := &bytes.Buffer{}
	err = ts.Execute(buf, app.addDefaultData(w, r, td))
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(chat.ContextUserKey).(models.User)
	if !ok {
		return nil
	}

	return &user
}

func (app *Application) deleteCookieAuthentication(w http.ResponseWriter, r *http.Request) {
	s, _ := app.Sessions.Get(r, sessionKey)
	delete(s.Values, "username")

	err := s.Save(r, w)
	if err != nil {
		app.serverError(w, err)
	}
}
