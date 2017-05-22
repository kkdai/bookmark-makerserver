package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
)

//BookmarkMgr :
type BookmarkMgr struct {
	User   string
	Repo   string
	Client *github.Client
}

//NewBookmark :
func NewBookmark(client *github.Client, user, repo string) *BookmarkMgr {
	new := new(BookmarkMgr)
	new.User = user
	new.Repo = repo
	new.Client = client
	return new
}

//CheckIfExist :
func (b *BookmarkMgr) CheckIfExist() bool {
	var issueByRepoOpts = &github.IssueListByRepoOptions{
		State:     "state",
		Direction: "asc",
	}

	issues, res, err := b.Client.Issues.ListByRepo(
		b.User,
		b.Repo,
		issueByRepoOpts,
	)

	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issues {
		fmt.Printf("#%3d: %s \n", *issue.Number, *issue.Title)
	}

	return false
}

//SaveBookmark :
func (b *BookmarkMgr) SaveBookmark() error {
	return nil
}
