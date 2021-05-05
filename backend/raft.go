package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"golang.org/x/net/context"
)

const hb = 1

type node struct {
	id     uint64
	ctx    context.Context
	store  *raft.MemoryStorage
	cfg    *raft.Config
	raft   raft.Node
	ticker <-chan time.Time
	done   chan struct{}

	articleList    map[string]*userArticleList
	articleRWMutex sync.RWMutex

	userList        []user
	userIDcount     int
	userListRWMutex sync.RWMutex

	follows        map[string]*userFollowList
	followsRWMutex sync.RWMutex
}

func newNode(id uint64, peers []raft.Peer) *node {
	store := raft.NewMemoryStorage()
	n := &node{
		id:    id,
		ctx:   context.TODO(),
		store: store,
		cfg: &raft.Config{
			ID:              id,
			ElectionTick:    10 * hb,
			HeartbeatTick:   hb,
			Storage:         store,
			MaxSizePerMsg:   math.MaxUint16,
			MaxInflightMsgs: 256,
		},
		ticker: time.Tick(time.Second),
		done:   make(chan struct{}),

		articleList: map[string]*userArticleList{
			"user1": &userArticleList{
				articleList: []article{
					article{ID: "1_1", User: "user1", PostDate: 1584525600, Title: "test article 1", Content: "this is the first test"},
				},
				articleID: 2,
				userID:    1,
			},
			"user2": &userArticleList{
				articleList: []article{
					article{ID: "2_1", User: "user2", PostDate: 1584439200, Title: "test article 2", Content: "this is the second test"},
				},
				articleID: 2,
				userID:    2,
			},
			"user3": &userArticleList{
				articleList: []article{
					article{ID: "3_1", User: "user3", PostDate: 1584352800, Title: "test article 3", Content: "this is the third test"},
				},
				articleID: 2,
				userID:    3,
			},
		},

		userList: []user{
			user{Username: "user1", Password: "pass1", UserID: 1},
			user{Username: "user2", Password: "pass2", UserID: 2},
			user{Username: "user3", Password: "pass3", UserID: 3},
		},

		userIDcount: 4,

		follows: map[string]*userFollowList{
			"user1": &userFollowList{followList: []string{"user2"}},
			"user2": &userFollowList{followList: []string{"user3"}},
			"user3": &userFollowList{followList: []string{"user1"}},
		},
	}

	n.raft = raft.StartNode(n.cfg, peers)
	return n
}

func (n *node) run() {
	for {
		select {
		case <-n.ticker:
			n.raft.Tick()
		case rd := <-n.raft.Ready():
			n.saveToStorage(rd.HardState, rd.Entries, rd.Snapshot)
			n.send(rd.Messages)
			if !raft.IsEmptySnap(rd.Snapshot) {
				n.processSnapshot(rd.Snapshot)
			}
			for _, entry := range rd.CommittedEntries {
				n.process(entry)
				if entry.Type == raftpb.EntryConfChange {
					var cc raftpb.ConfChange
					cc.Unmarshal(entry.Data)
					n.raft.ApplyConfChange(cc)
				}
			}
			n.raft.Advance()
		case <-n.done:
			return
		}
	}
}

func (n *node) saveToStorage(hardState raftpb.HardState, entries []raftpb.Entry, snapshot raftpb.Snapshot) {
	n.store.Append(entries)

	if !raft.IsEmptyHardState(hardState) {
		n.store.SetHardState(hardState)
	}

	if !raft.IsEmptySnap(snapshot) {
		n.store.ApplySnapshot(snapshot)
	}
}

func (n *node) send(messages []raftpb.Message) {
	for _, m := range messages {
		log.Println(raft.DescribeMessage(m, nil))

		// send message to other node
		nodes[int(m.To)].receive(n.ctx, m)
	}
}

func (n *node) processSnapshot(snapshot raftpb.Snapshot) {
	n.store.ApplySnapshot(snapshot)
}

func (n *node) process(entry raftpb.Entry) {
	log.Printf("node %v: processing entry: %v\n", n.id, entry)
	if entry.Type == raftpb.EntryNormal && entry.Data != nil {
		parts := bytes.SplitN(entry.Data, []byte(":"), 2)
		operation := string(parts[0])

		if operation == "addArticle" {

			var newArticle article
			json.Unmarshal(parts[1], &newArticle)

			title := newArticle.Title
			content := newArticle.Content
			user := newArticle.User
			timestampSeconds := newArticle.PostDate

			userAList := n.articleList[user]

			userAList.Lock()
			articleID := strconv.Itoa(userAList.userID) + "_" + strconv.Itoa(userAList.articleID)

			a := article{
				ID:       articleID,
				User:     user,
				PostDate: timestampSeconds,
				Title:    title,
				Content:  content,
			}

			userAList.articleID++
			userAList.articleList = append(userAList.articleList, a)
			userAList.Unlock()

		} else if operation == "addUser" {

			var newUser userProposal
			json.Unmarshal(parts[1], &newUser)

			username := newUser.Username
			password := newUser.Password

			n.userListRWMutex.Lock()
			u := user{Username: username, Password: password, UserID: n.userIDcount}
			n.userList = append(n.userList, u)
			n.userIDcount++
			n.userListRWMutex.Unlock()

			n.followsRWMutex.Lock()
			n.follows[username] = &userFollowList{}
			n.followsRWMutex.Unlock()

			n.articleRWMutex.Lock()
			n.articleList[username] = &userArticleList{articleID: 1, userID: u.UserID}
			n.articleRWMutex.Unlock()

		} else if operation == "addFollow" {

			var newFollow followAddProposal
			json.Unmarshal(parts[1], &newFollow)

			followUser := newFollow.FollowUser
			thisUser := newFollow.ThisUser

			f := n.follows[thisUser]

			f.Lock()
			f.followList = append(f.followList, followUser)
			f.Unlock()

		} else if operation == "removeFollow" {

			var removeFollow followRemoveProposal
			json.Unmarshal(parts[1], &removeFollow)

			thisUser := removeFollow.ThisUser
			i := removeFollow.RemoveIndex

			f := n.follows[thisUser]

			f.Lock()
			f.followList = append(f.followList[:i], f.followList[i+1:]...)
			f.Unlock()
		}
	}
}

func (n *node) receive(ctx context.Context, message raftpb.Message) {
	n.raft.Step(ctx, message)
}

var (
	nodes       = make(map[int]*node)
	currentNode = 3
	clusterSize = 3
)

func raftInit() {
	// start a small 3 node cluster
	nodes[1] = newNode(1, []raft.Peer{{ID: 1}, {ID: 2}, {ID: 3}})
	go nodes[1].run()

	nodes[2] = newNode(2, []raft.Peer{{ID: 1}, {ID: 2}, {ID: 3}})
	go nodes[2].run()

	nodes[3] = newNode(3, []raft.Peer{{ID: 1}, {ID: 2}, {ID: 3}})
	go nodes[3].run()

	// wait for cluster to set up
	time.Sleep(10 * time.Second)
	nodes[2].raft.Campaign(nodes[2].ctx)

	// demo code
	// go nodes[2].raft.Stop()
	// go nodes[1].raft.Stop()

	time.Sleep(10 * time.Second)
}

func getRaftNode() *node {
	time.Sleep(500 * time.Millisecond)

	// Do a check for one active Raft node?
	// Rotate between all active nodes
	currentNode = (currentNode % clusterSize) + 1
	i := 1
	for (nodes[currentNode].raft.Status().ID == 0) || (i < 3) {
		currentNode = (currentNode % clusterSize) + 1
		i++
	}
	fmt.Println("Using Raft Node#", currentNode)
	return nodes[currentNode]
}

func commitProposal(proposal []byte) error { // also include a channel in the arguments?
	// will keep retrying proposal until entry is committed
	committed := false
	var err error
	var node *node
	var oldCommitIndex, oldTerm uint64
	for !committed {

		node = getRaftNode()
		oldCommitIndex = node.raft.Status().Commit
		oldTerm = node.raft.Status().Term
		err = node.raft.Propose(node.ctx, proposal)

		for err != nil {
			time.Sleep(time.Second * 1)
			node = getRaftNode()
			err = node.raft.Propose(node.ctx, []byte(proposal))
			oldCommitIndex = node.raft.Status().Commit
			oldTerm = node.raft.Status().Term
		}

		for (node.raft.Status().Commit == oldCommitIndex) && (node.raft.Status().Term == oldTerm) {
			// waiting for node to commit or check if term changed
		}

		if (node.raft.Status().ID != 0) && (node.raft.Status().Term == oldTerm) {
			committed = true
		}
	}
	return nil
}

func foundEntry(entry []byte, n *node) bool {
	state := <-n.raft.Ready()
	for _, v := range state.CommittedEntries {
		if (bytes.Equal(v.Data, entry)) && (v.Type == 0) {
			return true
		}
	}
	return false
}
