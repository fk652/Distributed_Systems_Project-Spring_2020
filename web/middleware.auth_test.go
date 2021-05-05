package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoggedInUnauth(t *testing.T) {
	r := getRouter(false)
	r.GET("/", setLoggedIn(false), ensureLoggedIn(), func(c *gin.Context) {

		t.Fail()
	})

	testMiddlewareRequest(t, r, http.StatusUnauthorized)
}

func TestLoggedInAuth(t *testing.T) {
	r := getRouter(false)
	r.GET("/", setLoggedIn(true), ensureLoggedIn(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	testMiddlewareRequest(t, r, http.StatusOK)
}

func TestNotLoggedInAuth(t *testing.T) {
	r := getRouter(false)
	r.GET("/", setLoggedIn(true), ensureNotLoggedIn(), func(c *gin.Context) {
		t.Fail()
	})
	testMiddlewareRequest(t, r, http.StatusUnauthorized)
}

func TestNotLoggedInUnauth(t *testing.T) {
	r := getRouter(false)
	r.GET("/", setLoggedIn(false), ensureNotLoggedIn(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	testMiddlewareRequest(t, r, http.StatusOK)
}

func TestSetUserStatusAuth(t *testing.T) {
	r := getRouter(false)
	r.GET("/", setUserStatus(), func(c *gin.Context) {
		loggedInInterface, exists := c.Get("is_logged_in")
		if !exists || !loggedInInterface.(bool) {
			t.Fail()
		}
	})
	w := httptest.NewRecorder()

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "103"})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)
}

func TestSetUserStatusUnauth(t *testing.T) {
	r := getRouter(false)
	r.GET("/", setUserStatus(), func(c *gin.Context) {
		loggedInInterface, exists := c.Get("is_logged_in")
		if exists && loggedInInterface.(bool) {
			t.Fail()
		}
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(w, req)
}

func setLoggedIn(b bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("is_logged_in", b)
	}
}
