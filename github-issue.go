package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
)

//Bookmark :
type Bookmark struct {
	Name    string
	Link    string
	Comment string
}

//NewBookmark :
func NewBookmark(link string) *Bookmark {
	new := new(Bookmark)
	return new
}

//CheckIfExist :
func (b *Bookmark) CheckIfExist() bool {
	client := github.NewClient(nil)
	// ctx := context.Background()
	var issueByRepoOpts = &github.IssueListByRepoOptions{
		State:     "state",
		Direction: "asc",
	}

	issues, res, err := client.Issues.ListByRepo(
		"user",
		"repo",
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
func (b *Bookmark) SaveBookmark() error {
	return nil
}
