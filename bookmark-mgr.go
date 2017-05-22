package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gernest/mention"
	"github.com/google/go-github/github"
	"github.com/mvdan/xurls"
	"golang.org/x/oauth2"
)

func String(v string) *string { return &v }

//BookmarkMgr :
type BookmarkMgr struct {
	Token string
	User  string
	Repo  string
}

//NewBookmark :
func NewBookmark(user, repo, token string) *BookmarkMgr {
	new := new(BookmarkMgr)
	new.User = user
	new.Repo = repo
	new.Token = token
	return new
}

//CheckIfExist :
func (b *BookmarkMgr) CheckIfExist() bool {
	// var issueByRepoOpts = &github.IssueListByRepoOptions{
	// 	State:     "state",
	// 	Direction: "asc",
	// }

	// issues, res, err := b.Client.Issues.ListByRepo(
	// 	b.Client
	// 	b.User,
	// 	b.Repo,
	// 	issueByRepoOpts,
	// )

	// res.Body.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, issue := range issues {
	// 	fmt.Printf("#%3d: %s \n", *issue.Number, *issue.Title)
	// }

	return false
}

//SaveBookmark :
func (b *BookmarkMgr) SaveBookmark(tweet string) error {
	// token := os.Getenv("Token")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: b.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	links := xurls.Relaxed.FindAllString(tweet, -1)
	tags := mention.GetTags('#', strings.NewReader(tweet))
	var body string
	for _, v := range links {
		body = fmt.Sprintf("%s [link](%s)", body, v)
	}

	input := &github.IssueRequest{
		Title:    String(tweet),
		Body:     String(body),
		Assignee: String(""),
		Labels:   &tags,
	}

	issue, _, err := client.Issues.Create(ctx, b.User, b.Repo, input)
	if err != nil {
		fmt.Printf("Issues.Create returned error: %v", err)
		return err
	}

	fmt.Println(issue)
	return nil
}
