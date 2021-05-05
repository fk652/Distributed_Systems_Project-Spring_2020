package main

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/coreos/etcd/raft"
)

func raftInitTest() {
	if len(nodes) == 0 {
		raftInit()
	} else {
		nodes[1].raft = raft.RestartNode(nodes[1].cfg)
		nodes[2].raft = raft.RestartNode(nodes[2].cfg)
		nodes[3].raft = raft.RestartNode(nodes[3].cfg)
		time.Sleep(time.Second * 10)
	}
}

func createNewArticleTest(title string, content string, user string, timestampSeconds int64) error {

	ap := articleProposal{
		User:     user,
		PostDate: timestampSeconds,
		Title:    title,
		Content:  content,
	}
	bytes, _ := json.Marshal(ap)

	proposal := []byte("addArticle:" + string(bytes))
	err := commitProposal(proposal)
	if err != nil {
		return err
	}
	return nil
}

func Test1Node(t *testing.T) {
	raftInitTest()

	go nodes[1].raft.Stop()
	t.Log("Node 1 Killed")
	t.Log("one node killed")
	time.Sleep(5 * time.Second)

	preLength := len(nodes[3].articleList["user1"].articleList)

	title := "title"
	content := "content"
	user := "user1"
	timestampSeconds := time.Now().Unix()

	err := createNewArticleTest(title, content, user, timestampSeconds)

	if err != nil {
		t.Fail()
	}

	time.Sleep(5 * time.Second)
	postLength := len(nodes[3].articleList["user1"].articleList)

	if postLength == (preLength + 1) {
		t.Log("SUCCESS Raft functioning as intended \t preLength:", preLength, "postLength:", postLength)
	} else {
		t.Errorf("FAIL Raft not functioning %v %v", preLength, postLength)
	}

}

func TestConcurrent1Node(t *testing.T) {
	raftInitTest()

	preLength := len(nodes[3].articleList["user1"].articleList)
	var wg sync.WaitGroup

	for i := 1; i < 101; i++ {
		wg.Add(1)

		go func(i int) {
			if i == 50 {
				nodes[2].raft.Stop()
				t.Log("one node killed")
			}

			title := "title"
			content := "content"
			user := "user1"
			timestampSeconds := time.Now().Unix()

			err := createNewArticleTest(title, content, user, timestampSeconds)

			if err != nil {
				t.Errorf("%s FAIL proposing article %d", user, i)
			}

			defer wg.Done()
			time.Sleep(time.Second * 3)
		}(i)

	}
	wg.Wait()

	time.Sleep(5 * time.Second)
	postLength := len(nodes[3].articleList["user1"].articleList)

	if postLength > preLength {
		t.Log("SUCCESS Raft functioning as intended \t preLength:", preLength, "postLength:", postLength)
	} else {
		t.Errorf("FAIL Raft not functioning %v %v", preLength, postLength)
	}

	// For the demo
	// nodes[2].raft = raft.RestartNode(nodes[2].cfg)
	// preLength = len(nodes[2].articleList["user1"].articleList)
	// time.Sleep(10 * time.Second)
	// postLength = len(nodes[2].articleList["user1"].articleList)
	// t.Log("Node 2 preLength:", preLength, "postLength:", postLength)
}

func Test2Node(t *testing.T) {
	raftInitTest()

	go nodes[1].raft.Stop()
	t.Log("Node 1 Killed")
	go nodes[2].raft.Stop()
	t.Log("Node 2 Killed")
	t.Log("two nodes killed")
	time.Sleep(5 * time.Second)

	lastNode := 3
	preCommitLength := nodes[lastNode].raft.Status().Commit

	c1 := make(chan string, 1)
	go func() {
		c1 <- "result 1"

		title := "title"
		content := "content"
		user := "user1"
		timestampSeconds := time.Now().Unix()

		err := createNewArticleTest(title, content, user, timestampSeconds)

		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	time.Sleep(30 * time.Second)
	postCommitLength := nodes[lastNode].raft.Status().Commit

	if preCommitLength == postCommitLength {
		t.Log("SUCCESS Raft not functioning as intended \t preCommitLength:", preCommitLength, "postCommitLength:", postCommitLength)
	} else {
		t.Errorf("FAIL Raft still functioning \t preCommitLength: %v postCommitLength: %v", preCommitLength, postCommitLength)
	}
}
