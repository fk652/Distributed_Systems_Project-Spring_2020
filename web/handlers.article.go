// handlers.article.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pb "commonpb"
)

func showIndexPage(c *gin.Context) {

	render(c, gin.H{
		"title": "Index Page"}, "index.html")
}

func showWelcomePage(c *gin.Context, articles []*pb.Article) {

	render(c, gin.H{
		"title":     "Home Page",
		"home":      true,
		"localtime": localTimeConverter,
		"payload":   articles}, "home.html")

}

func showHomePage(c *gin.Context, articles []*pb.Article) {

	render(c, gin.H{
		"title":     "Home Page",
		"home":      true,
		"localtime": localTimeConverter,
		"payload":   articles}, "home.html")
}

func showHomePageEmpty(c *gin.Context) {

	render(c, gin.H{
		"title": "Home Page",
		"home":  true,
		"notify": []string{"Follow some people to see some interesting feed :)",
			"Try searching for users in the View Articles tab",
			"Think their tweets are interesting?",
			"Then smash that follow button!"},
	}, "home.html")
}

func showArticleCreationPage(c *gin.Context) {

	render(c, gin.H{
		"title":          "Create New Article",
		"article_create": true}, "create-article.html")
}

func showSearchArticlesPage(c *gin.Context) {

	render(c, gin.H{
		"title": "Search Articles",
		"view":  true,
	}, "search-articles.html")
}

func showSearchArticleError(c *gin.Context, err error) {

	c.HTML(http.StatusBadRequest, "search-articles.html", gin.H{
		"view":       true,
		"ErrorTitle": err.Error()})
}

func showArticle(c *gin.Context, a *pb.Article) {

	render(c, gin.H{
		"title":     a.Title,
		"view":      true,
		"localtime": localTimeConverter,
		"payload":   a}, "article.html")
}

func showArticleError(c *gin.Context, err error) {

	render(c, gin.H{
		"title":      "Article Error",
		"view":       true,
		"ErrorTitle": err.Error()}, "article.html")
}

func showUserArticle(c *gin.Context, username string, showFollowButton bool, followed bool, articles []*pb.Article) {

	if showFollowButton {

		render(c, gin.H{
			"title":            username + " Articles",
			"username":         username,
			"view":             true,
			"isFollowed":       followed,
			"showFollowButton": showFollowButton,
			"localtime":        localTimeConverter,
			"payload":          articles}, "user-articles.html")
	} else {
		render(c, gin.H{
			"title":      username + " Articles",
			"username":   username,
			"myArticles": true,
			"localtime":  localTimeConverter,
			"payload":    articles}, "user-articles.html")
	}
}

func showUserArticleError(c *gin.Context, err error) {

	c.HTML(http.StatusBadRequest, "user-articles.html", gin.H{
		"title":        "User Articles Error",
		"view":         true,
		"ErrorTitle":   "No Tweets Found",
		"ErrorMessage": err.Error()})
}

func showCreateArticleSuccess(c *gin.Context, a *pb.Article) {

	render(c, gin.H{
		"title":   "Submission Successful",
		"payload": a}, "submission-successful.html")
}
