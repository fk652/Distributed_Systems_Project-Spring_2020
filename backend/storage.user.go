// storage.user.go

package main

import "sync"

type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
	UserID   int    `json:"userID"`
}

type userProposal struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

var userList = []user{
	user{Username: "user1", Password: "pass1", UserID: 1},
	user{Username: "user2", Password: "pass2", UserID: 2},
	user{Username: "user3", Password: "pass3", UserID: 3},
}
var userIDcount = 4
var userListRWMutex sync.RWMutex

type userFollowList struct {
	sync.RWMutex
	followList []string
}

type followAddProposal struct {
	ThisUser   string `json:"thisuser"`
	FollowUser string `json:"followuser"`
}

type followRemoveProposal struct {
	ThisUser    string `json:"thisuser"`
	RemoveIndex int    `json:"removeindex"`
}

var follows = map[string]*userFollowList{
	"user1": &userFollowList{followList: []string{"user2"}},
	"user2": &userFollowList{followList: []string{"user3"}},
	"user3": &userFollowList{followList: []string{"user1"}},
}

var followsRWMutex sync.RWMutex
