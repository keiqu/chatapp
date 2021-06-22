package main

import "net/url"

type templateData struct {
	Errors map[string]string
	Form   url.Values
	Alert  string
}
