package forms

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testField = "test field"

func TestForm_ValidWhenNoErrors(t *testing.T) {
	f := New(url.Values{})

	assert.True(t, f.Valid())
}

func TestForm_ValidWhenThereIsErrors(t *testing.T) {
	f := New(url.Values{})
	f.Errors[testField] = "Invalid value."

	assert.False(t, f.Valid())
}

func TestForm_RequiredWhenFieldIsEmpty(t *testing.T) {
	uv := url.Values{}
	uv.Set(testField, "")
	f := New(uv)

	f.Required(testField)

	assert.NotEmpty(t, f.Errors[testField])
}

func TestForm_RequiredWhenFieldIsNotEmpty(t *testing.T) {
	uv := url.Values{}
	uv.Set(testField, "12345")
	f := New(uv)

	f.Required(testField)

	assert.Empty(t, f.Errors[testField])
}

func TestForm_MinLength(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		min         int
		isIncorrect bool
	}{
		{
			name:        "Empty string smaller than min",
			value:       "",
			min:         6,
			isIncorrect: true,
		},
		{
			name:        "Empty string equals min",
			value:       "",
			min:         0,
			isIncorrect: false,
		},
		{
			name:        "Length is smaller than min",
			value:       "123",
			min:         6,
			isIncorrect: true,
		},
		{
			name:        "Length is same as min",
			value:       "123456",
			min:         6,
			isIncorrect: false,
		},
		{
			name:        "Length is bigger than min",
			value:       "1234567",
			min:         6,
			isIncorrect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uv := url.Values{}
			uv.Set(testField, tt.value)
			f := New(uv)

			f.MinLength(testField, tt.min)

			_, ok := f.Errors[testField]
			assert.True(t, tt.isIncorrect == ok)
		})
	}
}

func TestForm_MaxLength(t *testing.T) {
	max := 6

	tests := []struct {
		name        string
		value       string
		isIncorrect bool
	}{
		{
			name:        "Empty string",
			value:       "",
			isIncorrect: false,
		},
		{
			name:        "Length is smaller than max",
			value:       "123",
			isIncorrect: false,
		},
		{
			name:        "Length is same as max",
			value:       "123456",
			isIncorrect: false,
		},
		{
			name:        "Length is bigger than max",
			value:       "1234567",
			isIncorrect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uv := url.Values{}
			uv.Set(testField, tt.value)
			f := New(uv)

			f.MaxLength(testField, max)

			_, ok := f.Errors[testField]
			assert.True(t, tt.isIncorrect == ok)
		})
	}
}

func TestForm_MatchesPattern(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		pattern     *regexp.Regexp
		isIncorrect bool
	}{
		{
			name:        "Empty string",
			value:       "",
			pattern:     EmailRX,
			isIncorrect: true,
		},
		{
			name:        "Matches pattern",
			value:       "username@example.com",
			pattern:     EmailRX,
			isIncorrect: false,
		},
		{
			name:        "Doesn't match pattern",
			value:       "123456",
			pattern:     EmailRX,
			isIncorrect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uv := url.Values{}
			uv.Set(testField, tt.value)
			f := New(uv)

			f.MatchesPattern(testField, tt.pattern)

			_, ok := f.Errors[testField]
			assert.True(t, tt.isIncorrect == ok)
		})
	}
}

func TestForm_ContainsOnlyAllowedChars(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		isIncorrect bool
	}{
		{
			name:        "Empty string",
			value:       "",
			isIncorrect: false,
		},
		{
			name:        "Contains only allowed characters",
			value:       "username123",
			isIncorrect: false,
		},
		{
			name:        "Contains >",
			value:       "username>",
			isIncorrect: true,
		},
		{
			name:        "Contains <",
			value:       "<username",
			isIncorrect: true,
		},
		{
			name:        "Contains &",
			value:       "user&name",
			isIncorrect: true,
		},
		{
			name:        "Contains '",
			value:       "usernam'e",
			isIncorrect: true,
		},
		{
			name:        "Contains \"",
			value:       "\"username",
			isIncorrect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uv := url.Values{}
			uv.Set(testField, tt.value)
			f := New(uv)

			f.ContainsOnlyAllowedChars(testField)

			_, ok := f.Errors[testField]
			assert.True(t, tt.isIncorrect == ok)
		})
	}
}
