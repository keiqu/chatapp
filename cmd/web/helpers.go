package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	log.Error().Msg(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
