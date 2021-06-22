package main

import "net/url"

type templateData struct {
	CSRFToken string
	Form      url.Values
	Alert     string
	Errors    map[string]string
}
