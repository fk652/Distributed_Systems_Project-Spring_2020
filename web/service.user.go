// service.user.go

package main

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	pb "github.com/fk652/import/commonpb"
)

func performLogin(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, _ := client.IsUserValid(ctx, &pb.AccountRequest{Username: username, Password: password})

	if r.GetReply() {

		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		c.SetCookie("username", username, 3600, "", "", false, true)

		showLoginSuccess(c)

	} else {

		showLoginFailed(c)
	}
}

func generateSessionToken() string {
	return strconv.FormatInt(rand.Int63(), 16)
}

func logout(c *gin.Context) {

	c.SetCookie("token", "", -1, "", "", false, true)
	c.SetCookie("username", "", -1, "", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func register(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, err := client.RegisterNewUser(ctx, &pb.AccountRequest{Username: username, Password: password})

	if r.GetReply() {

		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		c.SetCookie("username", username, 3600, "", "", false, true)

		showRegisterSuccess(c)

	} else {

		showRegisterFailed(c, err)
	}
}

func performFollow(c *gin.Context) {

	followUsername := c.PostForm("username")
	thisUser, _ := c.Cookie("username")
	// if err or thisUser == "" then return to home page
	if thisUser == followUsername {
		showFollowFailed(c, errors.New("Can't follow yourself"))
	}

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	_, err := client.AddFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})

	if err != nil {

		showFollowFailed(c, err)

	} else {

		showFollowSuccess(c, followUsername)
	}
}

func performUnfollow(c *gin.Context) {

	followUsername := c.PostForm("username")
	thisUser, _ := c.Cookie("username")
	// if err or thisUser == "" then return to home page
	if thisUser == followUsername {
		showUnfollowFailed(c, errors.New("Can't follow yourself"))
	}

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	_, err := client.RemoveFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})

	if err != nil {

		showUnfollowFailed(c, err)

	} else {

		showUnfollowSuccess(c, followUsername)
	}
}

func getFollowPage(c *gin.Context) {

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	thisUser, _ := c.Cookie("username")
	// if err or thisUser == "" then return to home page
	r, err := client.GetFollowedUsers(ctx, &pb.UsernameRequest{Username: thisUser})
	followedUsers := r.GetFollowList()

	if err != nil {

		showFollowPageError(c, err)

	} else if len(followedUsers) == 0 {

		showFollowPageEmpty(c)

	} else {

		showFollowPage(c, followedUsers)
	}
}
