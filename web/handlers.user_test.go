package main

import (
	//"io/ioutil"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestArticleCreatePage(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "111"})

	r.GET("/article/create", ensureLoggedIn(), showArticleCreationPage)

	req, _ := http.NewRequest("GET", "/article/create", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestArticleCreatePageUnauth(t *testing.T) {
	r := getRouter(true)

	r.GET("/article/create", ensureLoggedIn(), showArticleCreationPage)

	req, _ := http.NewRequest("GET", "/article/create", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusUnauthorized
	})
}

func TestArticleCreateUnauth(t *testing.T) {
	r := getRouter(true)

	r.POST("/article/create", ensureLoggedIn(), createArticle)

	articlePayload := getArticlePOSTPayload()
	req, _ := http.NewRequest("POST", "/article/create", strings.NewReader(articlePayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(articlePayload)))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusUnauthorized
	})
}

func getArticlePOSTPayload() string {
	params := url.Values{}
	params.Add("title", "Test Article Title")
	params.Add("content", "Test Article Content")

	return params.Encode()
}

func TestShowLoginPageAuth(t *testing.T) {
	w := httptest.NewRecorder()
	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.GET("/u/login", ensureNotLoggedIn(), showLoginPage)

	req, _ := http.NewRequest("GET", "/u/login", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

func TestShowLoginPageUnauth(t *testing.T) {
	r := getRouter(true)
	r.GET("/u/login", ensureNotLoggedIn(), showLoginPage)

	req, _ := http.NewRequest("GET", "/u/login", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Login</title>") > 0

		return statusOK && pageOK
	})
}

func TestLoginUnauthWrongID_Pass(t *testing.T) {
	w := httptest.NewRecorder()
	r := getRouter(true)

	r.POST("/u/login", ensureNotLoggedIn(), performLogin)

	loginPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/login", strings.NewReader(loginPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestRegisterAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "113"})
	r.POST("/u/register", ensureNotLoggedIn(), register)

	registrationPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

func TestRegisterUnauthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	r := getRouter(true)

	r.POST("/u/register", ensureNotLoggedIn(), register)

	registrationPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestRegisterUnauthUnavailableUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := getRouter(true)

	r.POST("/u/register", ensureNotLoggedIn(), register)

	registrationPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func getLoginPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "user1")
	params.Add("password", "pass1")

	return params.Encode()
}

func getRegistrationPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "u1")
	params.Add("password", "p1")

	return params.Encode()
}
