package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gernest/mention"
	"github.com/google/go-github/v34/github"
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
	return false
}

//SaveBookmark :
func (b *BookmarkMgr) SaveBookmark(tweet string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: b.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	links := xurls.Relaxed.FindAllString(tweet, -1)
	tags := mention.GetTags('#', strings.NewReader(tweet))
	title := fmt.Sprintf("%s", tweet)

	var body string
	var commentBody string
	strTs := strings.SplitN(tweet, "#", 2)

	if len(strTs) >= 2 {
		title = strTs[0]
		commentBody = strTs[1]
	}

	//To get pure comment, we need remove links and tags
	if commentBody != "" {
		for _, v := range links {
			commentBody = strings.Replace(commentBody, v, "", -1)
		}

		for _, v := range tags {
			commentBody = strings.Replace(commentBody, v, "", -1)
		}

		commentBody = strings.Replace(commentBody, "#", "", -1)
		commentBody = strings.TrimLeft(commentBody, " ")
	}

	//Prepare links, if no link just not post to github issue
	if len(links) == 0 {
		log.Println("Skip post:", tweet)
		return nil
	}

	for _, v := range links {
		body = fmt.Sprintf("%s [link](%s)", body, v)
	}

	//Add comment after links
	if commentBody != "" {
		body = fmt.Sprintf("%s \n %s", body, commentBody)
	}

	// Push to github issue
	if tags == nil {
		tags = []string{}
	}
	input := &github.IssueRequest{
		Title:    String(title),
		Body:     String(body),
		Assignee: String(""),
		Labels:   &tags,
	}

	_, _, err := client.Issues.Create(ctx, b.User, b.Repo, input)
	if err != nil {
		fmt.Printf("Issues.Create returned error: %v", err)
		return err
	}
	return nil
}
