package main

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	pb "github.com/fk652/import/commonpb"
)

func TestGetWelcomePage(t *testing.T) {
	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	r, err := client.GetAllArticles(ctx, &pb.Request{Message: "requesting all articles"})
	articles := r.GetArticles()
	if len(articles) == 0 {
		t.Fail()
	} else if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
}

func TestGetHomePage(t *testing.T) {
	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	thisUser := "user1"

	r, err := client.GetSomeArticles(ctx, &pb.UsernameRequest{Username: thisUser})
	articles := r.GetArticles()

	if len(articles) == 0 {
		t.Fail()
	} else if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
}

func TestGetArticle(t *testing.T) {

	articleID := "1_1"

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	_, err := client.GetArticleByID(ctx, &pb.ArticleIDRequest{Id: articleID})

	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
}

func TestGetUserArticle(t *testing.T) {

	username := "user1"

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, _ := client.IsUsernameAvailable(ctx, &pb.UsernameRequest{Username: username})
	check := r.GetReply()

	if !check {
		_, err := client.GetArticleByUser(ctx, &pb.UsernameRequest{Username: username})
		if err != nil {
			t.Fail()
		}
	}

}

func TestCreateArticle(t *testing.T) {
	title := "title"
	content := "content"

	thisUser := "user1"

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	timestamp := time.Now().Unix()
	_, err := client.CreateNewArticle(ctx, &pb.NewArticleRequest{Title: title, Content: content, User: thisUser, TimestampSeconds: timestamp})
	if err != nil {
		t.Fail()
	}

}

func TestGetAllArticles(t *testing.T) {

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, err := client.GetAllArticles(ctx, &pb.Request{Message: "requesting all articles"})
	articles := r.GetArticles()

	if len(articles) > 4 {
		t.Logf(strconv.Itoa(len(articles)))
		t.Fail()
	} else if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

}

func TestContextPassArticle(t *testing.T) {
	client, conn := connectToBackendServer()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := client.GetAllArticles(ctx, &pb.Request{Message: "requesting all articles"})
	if err != nil {
		t.Errorf("GetAllArticles FAIL")
	} else {
		t.Log("GetAllArticles PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	title := "context"
	content := "pass"
	thisUser := "user1"
	timestamp := time.Now().Unix()
	_, err = client.CreateNewArticle(ctx, &pb.NewArticleRequest{Title: title, Content: content, User: thisUser, TimestampSeconds: timestamp})
	if err != nil {
		t.Errorf("CreateNewArticle FAIL")
	} else {
		t.Log("CreateNewArticle PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	thisUser = "user3"
	_, err = client.GetSomeArticles(ctx, &pb.UsernameRequest{Username: thisUser})
	if err != nil {
		t.Errorf("GetSomeArticles FAIL")
	} else {
		t.Log("GetSomeArticles PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	articleID := "1_1"
	_, err = client.GetArticleByID(ctx, &pb.ArticleIDRequest{Id: articleID})
	if err != nil {
		t.Errorf("GetArticleByID FAIL")
	} else {
		t.Log("GetArticleByID PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	thisUser = "user2"
	_, err = client.GetArticleByUser(ctx, &pb.UsernameRequest{Username: thisUser})
	if err != nil {
		t.Errorf("GetArticleByUser FAIL")
	} else {
		t.Log("GetArticleByUser PASS")
	}
}

func TestContextFailArticle(t *testing.T) {
	client, conn := connectToBackendServer()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	_, err := client.GetAllArticles(ctx, &pb.Request{Message: "requesting all articles"})
	if err == nil {
		t.Errorf("GetAllArticles FAIL")
	} else {
		t.Log("GetAllArticles PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	title := "context"
	content := "fail"
	thisUser := "user1"
	timestamp := time.Now().Unix()
	_, err = client.CreateNewArticle(ctx, &pb.NewArticleRequest{Title: title, Content: content, User: thisUser, TimestampSeconds: timestamp})
	if err == nil {
		t.Errorf("CreateNewArticle FAIL")
	} else {
		t.Log("CreateNewArticle PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	thisUser = "user3"
	_, err = client.GetSomeArticles(ctx, &pb.UsernameRequest{Username: thisUser})
	if err == nil {
		t.Errorf("GetSomeArticles FAIL")
	} else {
		t.Log("GetSomeArticles PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	articleID := "1_1"
	_, err = client.GetArticleByID(ctx, &pb.ArticleIDRequest{Id: articleID})
	if err == nil {
		t.Errorf("GetArticleByID FAIL")
	} else {
		t.Log("GetArticleByID PASS")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
	cancel()
	thisUser = "user2"
	_, err = client.GetArticleByUser(ctx, &pb.UsernameRequest{Username: thisUser})
	if err == nil {
		t.Errorf("GetArticleByUser FAIL")
	} else {
		t.Log("GetArticleByUser PASS")
	}
}

func TestConcurrency(t *testing.T) {

	client, conn := connectToBackendServer()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	r, err := client.GetAllArticles(ctx, &pb.Request{Message: "requesting all articles"})
	oldLength := len(r.GetArticles())

	var wg sync.WaitGroup

	for i := 4; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			//Register User
			username := "username" + strconv.Itoa(i)
			password := "password" + strconv.Itoa(i)

			thisUser := username
			title := "title"
			content := "content"

			client, conn := connectToBackendServer()
			defer conn.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
			defer cancel()

			//register
			t.Log(username, "Registering new account")
			_, err := client.RegisterNewUser(ctx, &pb.AccountRequest{Username: username, Password: password})
			if err != nil {
				t.Errorf("%s FAIL Register", username)
			} else {
				t.Log(username, "Registered")
			}

			//Create Article
			timestamp := time.Now().Unix()
			for j := 0; j < 50; j++ {
				t.Log(username, "creating article", j)
				_, err := client.CreateNewArticle(ctx, &pb.NewArticleRequest{Title: title, Content: content, User: thisUser, TimestampSeconds: timestamp})
				if err != nil {
					t.Errorf("%s FAIL creating article %d", username, j)
				} else {
					t.Log(username, "created article", j)
				}
			}

			//Following
			followUsername := "user1"
			followUsername2 := "user2"
			followUsername3 := "user3"

			t.Log(username, "following user1")
			_, err1 := client.AddFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})
			if err1 != nil {
				t.Errorf("%s FAIL following user1", username)
			} else {
				t.Log(username, "followed user1")
			}

			t.Log(username, "following user2")
			_, err2 := client.AddFollow(ctx, &pb.FollowRequest{FollowUser: followUsername2, ThisUser: thisUser})
			if err2 != nil {
				t.Errorf("%s FAIL following user2", username)
			} else {
				t.Log(username, "followed user2")
			}

			t.Log(username, "following user3")
			_, err3 := client.AddFollow(ctx, &pb.FollowRequest{FollowUser: followUsername3, ThisUser: thisUser})
			if err3 != nil {
				t.Errorf("%s FAIL following user3", username)
			} else {
				t.Log(username, "followed user3")
			}

			t.Log(username, "unfollowing user1")
			_, err4 := client.RemoveFollow(ctx, &pb.FollowRequest{FollowUser: followUsername, ThisUser: thisUser})
			if err4 != nil {
				t.Errorf("%s FAIL unfollowing user1", username)
			} else {
				t.Log(username, "unfollowed user1")
			}

			defer wg.Done()
			time.Sleep(time.Second * 3)
		}(i)

	}
	wg.Wait()

	r, err = client.GetAllArticles(ctx, &pb.Request{Message: "requesting all articles"})
	newLength := len(r.GetArticles())
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	t.Log("Old Article Length:", oldLength, "\t New Article Length:", newLength)
}
