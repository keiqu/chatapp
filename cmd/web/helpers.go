package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"

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

func (app *application) render(w http.ResponseWriter, name string, td templateData) {
	ts, err := template.ParseFiles(filepath.Join("./ui/templates/", name), "./ui/templates/base.layout.gohtml")
	if err != nil {
		app.serverError(w, err)
	}

	buf := &bytes.Buffer{}
	err = ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		app.serverError(w, err)
	}
}
