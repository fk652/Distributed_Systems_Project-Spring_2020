// main.go

package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	pb "github.com/fk652/import/commonpb"
	"google.golang.org/grpc"
)

const (
	backendAddress = "localhost:50051"
	authAddress    = "localhost:50052"
)

var router *gin.Engine

func main() {

	//gin.SetMode(gin.ReleaseMode)

	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static("/img", "./img")

	initializeRoutes()

	router.Run()
}

func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":

		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":

		c.XML(http.StatusOK, data["payload"])
	default:

		c.HTML(http.StatusOK, templateName, data)
	}
}

func connectToBackendServer() (pb.BackendClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(backendAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewBackendClient(conn)

	return c, conn
}

func connectToAuthServer() (pb.AuthClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(authAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewAuthClient(conn)

	return c, conn
}
