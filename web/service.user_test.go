package main

import (
	"context"
	"math/rand"
	"testing"
	"time"

	pb "github.com/fk652/import/commonpb"
)

func TestPerformLogin(t *testing.T) {
	username := "user1"
	password := "pass1"

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, _ := client.IsUserValid(ctx, &pb.AccountRequest{Username: username, Password: password})

	if !r.GetReply() {
		t.Fail()
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestRegister(t *testing.T) {

	username := RandStringRunes(rand.Intn(9))
	password := RandStringRunes(rand.Intn(9))

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, err := client.RegisterNewUser(ctx, &pb.AccountRequest{Username: username, Password: password})

	if !r.GetReply() {
		t.Logf(err.Error())
		t.Fail()
	}
}

func TestPerformFollow(t *testing.T) {

	followUsername := "user3"

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	thisUser := "user1"

	_, err := client.AddFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})

	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

}

func TestPerformUnFollow(t *testing.T) {
	followUsername := "user3"

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	thisUser := "user1"

	_, err := client.RemoveFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})

	if err != nil {
		t.Fail()
	}
}

func TestGetFollowPage(t *testing.T) {

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	thisUser := "user1"

	_, err := client.GetFollowedUsers(ctx, &pb.UsernameRequest{Username: thisUser})

	if err != nil {
		t.Fail()
	}

}

func TestContextFailUser(t *testing.T) {
	client, conn := connectToBackendServer()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	thisUser := "user1"
	_, err := client.GetFollowedUsers(ctx, &pb.UsernameRequest{Username: thisUser})
	if err == nil {
		t.Errorf("GetFollowedUsers FAIL")
	} else {
		t.Log("GetFollowedUsers PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	followUsername := "user3"
	thisUser = "user1"
	_, err = client.RemoveFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})
	if err == nil {
		t.Errorf("RemoveFollow FAIL")
	} else {
		t.Log("RemoveFollow PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	followUsername = "user3"
	thisUser = "user2"
	_, err = client.AddFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})
	if err == nil {
		t.Errorf("AddFollow FAIL")
	} else {
		t.Log("AddFollow PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	username := RandStringRunes(rand.Intn(9))
	password := RandStringRunes(rand.Intn(9))
	_, err = client.RegisterNewUser(ctx, &pb.AccountRequest{Username: username, Password: password})
	if err == nil {
		t.Errorf("RegisterNewUser FAIL")
	} else {
		t.Log("RegisterNewUser PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	username = "user1"
	password = "pass1"
	_, err = client.IsUserValid(ctx, &pb.AccountRequest{Username: username, Password: password})
	if err == nil {
		t.Errorf("IsUserValid FAIL")
	} else {
		t.Log("IsUserValid PASS")
	}
}