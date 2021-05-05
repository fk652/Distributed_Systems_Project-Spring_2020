// middleware.auth.go

package main

import (
	"context"
	"net/http"
	"time"

	pb "github.com/fk652/import/commonpb"
	"github.com/gin-gonic/gin"
)

type contextArg struct {
	c *gin.Context
}

func ensureLoggedIn() gin.HandlerFunc {

	return func(c *gin.Context) {

		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)

		client, conn := connectToAuthServer()
		defer conn.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		r, err := client.EnsureLoggedIn(ctx, &pb.BoolRequest{Request: loggedIn})

		if err != nil || !r.Reply {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func ensureNotLoggedIn() gin.HandlerFunc {

	return func(c *gin.Context) {

		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)

		client, conn := connectToAuthServer()
		defer conn.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		r, err := client.EnsureNotLoggedIn(ctx, &pb.BoolRequest{Request: loggedIn})

		if err != nil || !r.Reply {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func setUserStatus() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := c.Cookie("token")

		if err != nil {
			c.Set("is_logged_in", false)
		} else {

			client, conn := connectToAuthServer()
			defer conn.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
			defer cancel()
			r, err := client.SetUserStatus(ctx, &pb.Request{Message: token})

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else if r.Reply {
				c.Set("is_logged_in", true)
			} else {
				c.Set("is_logged_in", false)
			}
		}
	}
}
