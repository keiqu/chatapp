package server

import (
	"html"
	"net/http"
	"net/url"
	"testing"

	"github.com/lazy-void/chatapp/models/mock"

	"github.com/stretchr/testify/assert"
)

func TestApplication_Home(t *testing.T) {
	app := newTestApp()

	tests := []struct {
		name            string
		isAuthenticated bool
		wantCode        int
		wantLocation    string
	}{
		{"Non-authenticated user", false, http.StatusSeeOther, "/user/login"},
		{"Authenticated user", true, http.StatusOK, ""},
	}

	for _, tt := range tests {
		tt := tt // create new variable for each closure

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := newTestServer(t, app.NewRouter())
			if tt.isAuthenticated {
				ts.authenticate(t)
			}

			code, header, _ := ts.get(t, "/")

			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantLocation, header.Get("Location"))
		})
	}
}

func TestApplication_SignupUser(t *testing.T) {
	app := newTestApp()

	tests := []struct {
		name                      string
		username, email, password string
		wantCode                  int
		wantBody                  string
	}{
		{
			"Valid submission",
			"Alice",
			"alice@gmail.com",
			"validPass123",
			http.StatusSeeOther,
			"",
		},
		{
			"Empty username",
			"",
			"alice@gmail.com",
			"validPass123",
			http.StatusOK,
			"This field cannot be empty.",
		},
		{
			"Empty email",
			"Alice",
			"",
			"validPass123",
			http.StatusOK,
			"This field cannot be empty.",
		},
		{
			"Empty password",
			"Alice",
			"alice@gmail.com",
			"",
			http.StatusOK,
			"This field cannot be empty.",
		},
		{
			"Username contains js",
			"<script>alert('hello');</script>",
			"alice@gmail.com",
			"validPass123",
			http.StatusOK,
			"Characters <, >, &, ' and \" are not allowed.",
		},
		{
			"Username contains html",
			"<strong>Alice</strong>",
			"alice@gmail.com",
			"validPass123",
			http.StatusOK,
			"Characters <, >, &, ' and \" are not allowed.",
		},
		{
			"Username is too long",
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			"alice@gmail.com",
			"validPass123",
			http.StatusOK,
			"Provided value is too long (maximum is 50 characters)",
		},
		{
			"Email is invalid (incomplete domain)",
			"Alice",
			"alice@gmail.",
			"validPass123",
			http.StatusOK,
			"Incorrect value.",
		},
		{
			"Email is invalid (missing @)",
			"Alice",
			"alicegmail.com",
			"validPass123",
			http.StatusOK,
			"Incorrect value.",
		},
		{
			"Email is invalid (missing username)",
			"Alice",
			"@gmail.com",
			"validPass123",
			http.StatusOK,
			"Incorrect value.",
		},
		{
			"Password is too short",
			"Alice",
			"alice@gmail.com",
			"123456789",
			http.StatusOK,
			"Provided value is too short (minimum is 10 characters)",
		},
		{
			"Duplicate username",
			mock.DupeUsername,
			"alice@gmail.com",
			"validPass123",
			http.StatusOK,
			"Username is already taken.",
		},
		{
			"Duplicate email",
			"Alice",
			mock.DupeEmail,
			"validPass123",
			http.StatusOK,
			"Email is already in use.",
		},
	}

	for _, tt := range tests {
		tt := tt // create new variable for each closure

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := newTestServer(t, app.NewRouter())

			_, _, body := ts.get(t, "/user/signup")
			csrfToken, err := extractCSRFToken(body)
			if err != nil {
				t.Fatal(err)
			}

			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", csrfToken)

			code, _, body := ts.post(t, "/user/signup", form)

			assert.Equal(t, tt.wantCode, code)
			assert.Contains(t, html.UnescapeString(body), tt.wantBody)
		})
	}
}

func TestApplication_LoginUser(t *testing.T) {
	app := newTestApp()

	tests := []struct {
		name            string
		email, password string
		wantCode        int
		wantBody        string
	}{
		{
			"Valid submission",
			mock.UserMock.Email,
			mock.ValidPassword,
			http.StatusSeeOther,
			"",
		},
		{
			"Empty email",
			"",
			mock.ValidPassword,
			http.StatusOK,
			"This field cannot be empty.",
		},
		{
			"Empty password",
			mock.UserMock.Email,
			"",
			http.StatusOK,
			"This field cannot be empty.",
		},
		{
			"Non-existent email",
			"alice@gmail.com",
			mock.ValidPassword,
			http.StatusOK,
			"Account with such email doesn't exist.",
		},
		{
			"Incorrect password",
			mock.UserMock.Email,
			"invalidPass123",
			http.StatusOK,
			"Incorrect password.",
		},
	}

	for _, tt := range tests {
		tt := tt // create new variable for each closure

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := newTestServer(t, app.NewRouter())

			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			_, _, body := ts.get(t, "/user/login")
			csrfToken, err := extractCSRFToken(body)
			if err != nil {
				t.Fatal(err)
			}

			form.Add("csrf_token", csrfToken)

			code, _, body := ts.post(t, "/user/login", form)

			assert.Equal(t, tt.wantCode, code)
			assert.Contains(t, html.UnescapeString(body), tt.wantBody)
		})
	}
}

func TestApplication_LogoutUser(t *testing.T) {
	app := newTestApp()

	tests := []struct {
		name            string
		isAuthenticated bool
		wantCode        int
		wantLocation    string
	}{
		{
			"Authenticated user",
			true,
			http.StatusSeeOther,
			"/user/login",
		},
		{
			"Non-authenticated user (no csrf token)",
			false,
			http.StatusBadRequest,
			"",
		},
	}

	for _, tt := range tests {
		tt := tt // create new variable for each closure

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := newTestServer(t, app.NewRouter())

			form := url.Values{}
			if tt.isAuthenticated {
				ts.authenticate(t)

				// Collect csrf token and cookie
				_, _, body := ts.get(t, "/")
				csrfToken, err := extractCSRFToken(body)
				if err != nil {
					t.Fatal(err)
				}
				form.Add("csrf_token", csrfToken)
			}

			code, header, _ := ts.post(t, "/user/logout", form)

			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantLocation, header.Get("Location"))
		})
	}
}

func TestApplication_LoginUserFormShowsSuccessFlashAfterRegistration(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	ts := newTestServer(t, app.NewRouter())

	// signup
	_, _, body := ts.get(t, "/user/signup")
	csrfToken, err := extractCSRFToken(body)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("username", "Alice")
	form.Add("email", "alice@gmail.com")
	form.Add("password", "validPassword123")
	form.Add("csrf_token", csrfToken)

	ts.post(t, "/user/signup", form)

	// get login page
	_, _, body = ts.get(t, "/user/login")

	// check success flash
	assert.Contains(t, body, "Your signup was successful. Please log in.")
}

func TestApplicationWhenNoCSRFTokenFormsRespondWithBadRequest(t *testing.T) {
	app := newTestApp()

	tests := []struct {
		name string
		path string
		form url.Values
	}{
		{
			"Signup form",
			"/user/signup",
			url.Values{
				"username": []string{"alice"},
				"email":    []string{"alice@gmail.com"},
				"password": []string{"validPassword123"},
			},
		},
		{
			"Login form",
			"/user/login",
			url.Values{
				"email":    []string{mock.UserMock.Email},
				"password": []string{mock.ValidPassword},
			},
		},
		{
			"Logout form",
			"/user/signup",
			url.Values{},
		},
	}

	for _, tt := range tests {
		tt := tt // create new variable for each closure

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := newTestServer(t, app.NewRouter())

			code, _, _ := ts.post(t, "/user/login", tt.form)

			assert.Equal(t, http.StatusBadRequest, code)
		})
	}
}

func TestApplicationWhenAuthenticatedUserRequestsPathThatRequiresToBeNonAuthenticated(t *testing.T) {
	app := newTestApp()

	// we do not test POST methods because they
	// require CSRF token from the GET counterparts
	tests := []struct {
		name   string
		path   string
		method string
	}{
		{"GET /user/signup", "/user/signup", http.MethodGet},
		{"GET /user/login", "/user/login", http.MethodGet},
	}

	for _, tt := range tests {
		tt := tt // create new variable for each closure

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := newTestServer(t, app.NewRouter())
			ts.authenticate(t)

			code, header, _ := ts.get(t, tt.path)

			assert.Equal(t, http.StatusSeeOther, code)
			assert.Equal(t, "/", header.Get("Location"))
		})
	}
}
