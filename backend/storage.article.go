// storage.article.go

package main

import "sync"

type article struct {
	ID       string `json:"id"`
	User     string `json:"user"`
	PostDate int64  `json:"postdate"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type articleProposal struct {
	User     string `json:"user"`
	PostDate int64  `json:"postdate"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type userArticleList struct {
	sync.RWMutex
	articleList []article
	articleID   int
	userID      int
}

var articleList = map[string]*userArticleList{
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
}

var articleRWMutex sync.RWMutex
