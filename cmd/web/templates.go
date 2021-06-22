package main

import "net/url"

type templateData struct {
	CSRFToken    string
	Form         url.Values
	SuccessFlash string
	ErrorFlash   string
	Errors       map[string]string
}
