package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
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

func (s *Server) createMessages(in <-chan string) {
	for {
		msg := <-in

		_, err := s.Messages.Insert(msg, time.Now().UTC())
		if err != nil {
			log.Err(err)
		}
	}
}
