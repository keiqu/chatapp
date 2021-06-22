package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/justinas/nosurf"

	"github.com/rs/zerolog/log"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	log.Error().Msg(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (app *application) addDefaultData(w http.ResponseWriter, r *http.Request, td templateData) templateData {
	td.CSRFToken = nosurf.Token(r)

	s, _ := app.sessions.Get(r, "user")
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

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td templateData) {
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

func (app *application) isAuthenticated(r *http.Request) bool {
	s, err := app.sessions.Get(r, "user")
	if err != nil {
		return false
	}

	_, ok := s.Values["userID"]
	return ok
}
