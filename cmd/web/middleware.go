package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func noSurf(next http.Handler) http.Handler {
	return nosurf.New(next)
}
