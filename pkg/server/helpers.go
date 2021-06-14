package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (s *Server) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	s.Error(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *Server) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
