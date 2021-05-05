// handlers.user.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func showLoginPage(c *gin.Context) {

	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func showLoginSuccess(c *gin.Context) {

	render(c, gin.H{
		"title": "Successful Login",
	}, "login-successful.html")
}

func showLoginFailed(c *gin.Context) {

	c.HTML(http.StatusBadRequest, "login.html", gin.H{
		"ErrorTitle":   "Login Failed",
		"ErrorMessage": "Invalid credentials provided"})
}

func showRegistrationPage(c *gin.Context) {

	render(c, gin.H{
		"title": "Register"}, "register.html")
}

func showRegisterSuccess(c *gin.Context) {

	render(c, gin.H{
		"title": "Successful registration & Login",
	}, "registration-successful.html")
}

func showRegisterFailed(c *gin.Context, err error) {

	c.HTML(http.StatusBadRequest, "register.html", gin.H{
		"ErrorTitle":   "Registration Failed",
		"ErrorMessage": err.Error()})
}

func showFollowSuccess(c *gin.Context, username string) {

	render(c, gin.H{
		"title":         "Follow Submit",
		"FollowSuccess": true,
		"FollowedUser":  username,
		"followed":      true}, "follow-submission.html")

}

func showFollowFailed(c *gin.Context, err error) {

	render(c, gin.H{
		"title":        "Follow Submit",
		"ErrorTitle":   "Follow Failed",
		"ErrorMessage": err.Error(),
		"followed":     true}, "follow-submission.html")
}

func showUnfollowSuccess(c *gin.Context, username string) {

	render(c, gin.H{
		"title":           "Unfollow Submit",
		"UnfollowSuccess": true,
		"UnfollowedUser":  username,
		"followed":        true}, "follow-submission.html")
}

func showUnfollowFailed(c *gin.Context, err error) {

	render(c, gin.H{
		"title":        "Unfollow Submit",
		"ErrorTitle":   "Unfollow Failed",
		"ErrorMessage": err.Error(),
		"followed":     true}, "follow-submission.html")
}

func showFollowPage(c *gin.Context, followedUsers []string) {

	render(c, gin.H{
		"title":    "Followed",
		"payload":  followedUsers,
		"followed": true}, "followed.html")
}

func showFollowPageEmpty(c *gin.Context) {

	render(c, gin.H{
		"title":  "Followed",
		"notify": "Your following list is empty",
		// "payload":  [],
		"followed": true}, "followed.html")
}

func showFollowPageError(c *gin.Context, err error) {

	render(c, gin.H{
		"title":        "Followed",
		"ErrorTitle":   "Re Failed",
		"ErrorMessage": err.Error(),
		// "payload":  [],
		"followed": true}, "followed.html")
}
