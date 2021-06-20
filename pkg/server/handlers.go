package server

import (
	"html/template"
	"net/http"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/templates/home.page.gohtml", "./ui/templates/base.layout.gohtml")
	if err != nil {
		s.serverError(w, err)
		return
	}
	messages, err := s.Messages.Get(100, 0)
	if err != nil {
		s.serverError(w, err)
		return
	}

	err = ts.Execute(w, &templateData{messages})
	if err != nil {
		s.serverError(w, err)
	}
}
