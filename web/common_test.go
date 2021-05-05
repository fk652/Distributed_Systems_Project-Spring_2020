package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

type article struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	PostDate string `json:"postdate"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

var articleList = []article{}

var userList = []user{
	user{Username: "user1", Password: "pass1"},
	user{Username: "user2", Password: "pass2"},
	user{Username: "user3", Password: "pass3"},
}

var tmpUserList []user
var tmpArticleList []article

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("templates/*")
		r.Use(setUserStatus())
	}
	return r
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func testMiddlewareRequest(t *testing.T, r *gin.Engine, expectedHTTPCode int) {
	req, _ := http.NewRequest("GET", "/", nil)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == expectedHTTPCode
	})
}

func saveLists() {
	tmpUserList = userList
	tmpArticleList = articleList
}

func restoreLists() {
	userList = tmpUserList
	articleList = tmpArticleList
}
