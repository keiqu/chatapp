package server

import (
	"errors"
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/lazy-void/chatapp/models/mock"
)

func extractCSRFToken(body string) (string, error) {
	regex := regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.*)">`)
	matches := regex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("csrf token was not found")
	}

	return html.UnescapeString(matches[1]), nil
}

func newTestApp() *Application {
	return &Application{
		Sessions: sessions.NewCookieStore([]byte("946IpCV9y5Vlur8YvODJEhaOY8m9J1E4")),
		Messages: &mock.MessageModel{},
		Users:    &mock.UserModel{},
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, router http.Handler) testServer {
	ts := httptest.NewServer(router)

	// store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar

	// no redirects
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return testServer{ts}
}

func (ts testServer) get(t *testing.T, path string) (int, http.Header, string) {
	resp, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return resp.StatusCode, resp.Header, string(body)
}

func (ts testServer) post(t *testing.T, path string, form url.Values) (int, http.Header, string) {
	resp, err := ts.Client().PostForm(ts.URL+path, form)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return resp.StatusCode, resp.Header, string(body)
}

func (ts testServer) authenticate(t *testing.T) {
	_, _, body := ts.get(t, "/user/login")
	csrfToken, err := extractCSRFToken(body)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("email", mock.UserMock.Email)
	form.Add("password", mock.ValidPassword)
	form.Add("csrf_token", csrfToken)

	_, _, _ = ts.post(t, "/user/login", form)

	// we need to delete csrf_token cookie that was set on login
	// because it may interfere with following requests that require
	// csrf token and use path that starts with /user
	u, err := url.Parse(ts.URL + "/user")
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar.SetCookies(u, []*http.Cookie{{
		Name:    "csrf_token",
		Value:   "",
		Path:    "/user",
		Expires: time.Unix(0, 0),
	}})

}
