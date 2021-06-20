package main

import (
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/templates/home.page.gohtml", "./ui/templates/base.layout.gohtml")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
