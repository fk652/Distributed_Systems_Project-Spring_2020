// service.article.go

package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"context"

	pb "github.com/fk652/import/commonpb"
)

func getWelcomePage(c *gin.Context) {

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, _ := client.GetAllArticles(ctx, &pb.Request{Message: "requesting all articles"})
	articles := r.GetArticles()

	showWelcomePage(c, articles)
}

func getHomePage(c *gin.Context) {

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	thisUser, _ := c.Cookie("username")
	// if err or thisUser == "" then return to home page

	r, _ := client.GetSomeArticles(ctx, &pb.UsernameRequest{Username: thisUser})
	articles := r.GetArticles()

	if len(articles) == 0 {

		showHomePageEmpty(c)

	} else {

		showHomePage(c, articles)
	}
}

func performSearchArticles(c *gin.Context) {

	username := c.PostForm("username")
	articleID := c.PostForm("articleID")

	if articleID != "" {
		location := "/article/articleID/" + articleID
		c.Redirect(http.StatusTemporaryRedirect, location)

	} else if username != "" {

		location := "/article/userArticles/" + username
		c.Redirect(http.StatusTemporaryRedirect, location)

	} else {

		showSearchArticleError(c, errors.New("No input found"))
	}
}

func getArticle(c *gin.Context) {

	articleID := c.Param("article_id")

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, err := client.GetArticleByID(ctx, &pb.ArticleIDRequest{Id: articleID})
	a := r.GetArticle()

	if err == nil {

		showArticle(c, a)

	} else {

		showArticleError(c, errors.New("Article not found"))
	}
}

func getUserArticle(c *gin.Context) {

	username := c.Param("username")

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, _ := client.IsUsernameAvailable(ctx, &pb.UsernameRequest{Username: username})
	check := r.GetReply()

	if !check {

		r, err := client.GetArticleByUser(ctx, &pb.UsernameRequest{Username: username})
		articles := r.GetArticles()

		if len(articles) > -1 {

			thisUser, _ := c.Cookie("username")
			// if err or thisUser == "" then return to home page

			r, _ := client.IsFollowed(ctx, &pb.FollowRequest{FollowUser: username, ThisUser: thisUser})
			followed := r.GetFound()

			showFollowButton := thisUser != username

			showUserArticle(c, username, showFollowButton, followed, articles)

		} else {

			showUserArticleError(c, err)
		}

	} else {

		showUserArticleError(c, errors.New("This user does not exist"))
	}
}

func getMyArticles(c *gin.Context) {

	thisUser, err := c.Cookie("username")

	if err != nil {

		c.AbortWithStatus(http.StatusNotFound)

	} else {

		location := "/article/userArticles/" + thisUser
		c.Redirect(http.StatusTemporaryRedirect, location)
	}
}

func createArticle(c *gin.Context) {

	title := c.PostForm("title")
	content := c.PostForm("content")

	thisUser, userErr := c.Cookie("username")

	if userErr == nil {

		client, conn := connectToBackendServer()
		defer conn.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		timestamp := time.Now().Unix()
		r, err := client.CreateNewArticle(ctx, &pb.NewArticleRequest{Title: title, Content: content, User: thisUser, TimestampSeconds: timestamp})

		if err == nil {

			a := r.GetArticle()
			showCreateArticleSuccess(c, a)

		} else {

			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {

		c.AbortWithStatus(http.StatusBadRequest)
	}
}

// timestamp converter
type localTime struct {
	Convert func(int64) string
}

var localTimeConverter = localTime{
	Convert: func(timestamp int64) string {
		return time.Unix(timestamp, 0).Format("Mon Jan 2 2006 15:04 PM")
		// timestamp := time.Now().Format("Mon Jan 2 2006 15:04 PM")
	},
}
