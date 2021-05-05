// middleware.auth.go

package main

import (
	"context"
	"fmt"

	pb "github.com/fk652/import/commonpb"
	"github.com/gin-gonic/gin"
)

type contextArg struct {
	c *gin.Context
}

func (s *server) EnsureLoggedIn(ctx context.Context, args *pb.BoolRequest) (*pb.BoolReply, error) {
	fmt.Println("EnsureLoggedIn")

	loggedIn := args.GetRequest()

	if !loggedIn {
		return &pb.BoolReply{Reply: false}, nil
	}
	return &pb.BoolReply{Reply: true}, nil
}

func (s *server) EnsureNotLoggedIn(ctx context.Context, args *pb.BoolRequest) (*pb.BoolReply, error) {
	fmt.Println("EnsureNotLoggedIn")

	loggedIn := args.GetRequest()

	if loggedIn {
		return &pb.BoolReply{Reply: false}, nil
	}
	return &pb.BoolReply{Reply: true}, nil
}

func (s *server) SetUserStatus(ctx context.Context, args *pb.Request) (*pb.BoolReply, error) {
	fmt.Println("SetUserStatus")

	token := args.GetMessage()

	if token != "" {
		return &pb.BoolReply{Reply: true}, nil
	}
	return &pb.BoolReply{Reply: false}, nil
}
