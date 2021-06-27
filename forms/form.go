// Package forms implements validation for POST forms.
package forms

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"unicode/utf8"
)

type errors map[string]string

func (e errors) add(field, error string) {
	// do not add error if field already has one
	if _, ok := e[field]; ok {
		return
	}

	e[field] = error
}

// EmailRX is regular expression for checking correctness of the email.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form wraps around postForm and allows to
// test its fields.
type Form struct {
	url.Values

	// Contains errors for each field.
	Errors errors
}

// New initializes Form.
func New(postForm url.Values) Form {
	return Form{
		Values: postForm,
		Errors: make(map[string]string),
	}
}

// Required adds an error if the field is empty.
func (f Form) Required(field string) {
	if v := f.Get(field); v != "" {
		return
	}

	f.Errors.add(field, "This field cannot be empty.")
}

// MinLength adds an error if the field contains
// less characters than min.
func (f Form) MinLength(field string, min int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) >= min {
		return
	}

	f.Errors.add(field, fmt.Sprintf("Provided value is too short (minimum is %d characters)", min))
}

// MaxLength adds an error if the field contains
// more characters than min.
func (f Form) MaxLength(field string, max int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) <= max {
		return
	}

	f.Errors.add(field, fmt.Sprintf("Provided value is too long (maximum is %d characters)", max))
}

// MatchesPattern adds an error if the field doesn't match
// the provided pattern.
func (f Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if pattern.MatchString(value) {
		return
	}

	f.Errors.add(field, "Incorrect value.")
}

// ContainsOnlyAllowedChars adds an error if the field contains one of
// the characters: <, >, &, ' and ". Allows to test for html/js embedding.
func (f Form) ContainsOnlyAllowedChars(field string) {
	value := f.Get(field)
	if html.EscapeString(value) == value {
		return
	}

	f.Errors.add(field, "Characters <, >, &, ' and \" are not allowed.")
}

// Valid returns true if form doesn't have
// any errors and false otherwise.
func (f Form) Valid() bool {
	return len(f.Errors) == 0
}
