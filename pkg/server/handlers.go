package server

import (
	"html/template"
	"net/http"
	"time"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/templates/home.page.gohtml", "./ui/templates/base.layout.gohtml")
	if err != nil {
		s.serverError(w, err)
		return
	}
	messages, err := s.Messages.Latest(100)
	if err != nil {
		s.serverError(w, err)
		return
	}

	err = ts.Execute(w, &templateData{messages})
	if err != nil {
		s.serverError(w, err)
	}
}

func (s *Server) createMessage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	msg := r.PostForm.Get("message")
	if msg == "" {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	_, err = s.Messages.Insert(msg, time.Now().UTC())
	if err != nil {
		s.serverError(w, err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
