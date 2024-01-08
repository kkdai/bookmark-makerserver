package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gernest/mention"
	"github.com/google/go-github/v57/github"
	"github.com/mvdan/xurls"
	"golang.org/x/oauth2"
)

// BookmarkMgr manages bookmarks for a GitHub repository.
type BookmarkMgr struct {
	Token string
	User  string
	Repo  string
}

// NewBookmark creates a new BookmarkMgr instance.
func NewBookmark(user, repo, token string) *BookmarkMgr {
	return &BookmarkMgr{
		User:  user,
		Repo:  repo,
		Token: token,
	}
}

// SaveBookmark saves a bookmark as a GitHub issue.
func (b *BookmarkMgr) SaveBookmark(tweet string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: b.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	links := xurls.Relaxed.FindAllString(tweet, -1)
	tags := mention.GetTags('#', strings.NewReader(tweet))
	title := tweet

	if strings.Contains(tweet, "#") {
		title = strings.SplitN(tweet, "#", 2)[0]
	}

	// Prepare the body of the issue by removing links and tags from the comment.
	commentBody := strings.Replace(tweet, title, "", 1)
	for _, v := range links {
		commentBody = strings.Replace(commentBody, v, "", -1)
	}
	for _, v := range tags {
		commentBody = strings.Replace(commentBody, v, "", -1)
	}
	commentBody = strings.TrimSpace(strings.Replace(commentBody, "#", "", -1))

	// If no link is present, skip posting to GitHub.
	if len(links) == 0 {
		log.Printf("Skip post: %s", tweet)
		return nil
	}

	var bodyBuilder strings.Builder
	for _, link := range links {
		bodyBuilder.WriteString(fmt.Sprintf(" [link](%s)", link))
	}
	if commentBody != "" {
		bodyBuilder.WriteString(fmt.Sprintf("\n%s", commentBody))
	}

	// Check tags if nil, apply default tag with #twitter hashtag.
	if len(tags) == 0 {
		tags = []string{"twitter"}
	}

	// Create a GitHub issue.
	input := &github.IssueRequest{
		Title:  &title,
		Body:   github.String(bodyBuilder.String()),
		Labels: &tags,
	}

	_, _, err := client.Issues.Create(ctx, b.User, b.Repo, input)
	if err != nil {
		log.Printf("Issues.Create returned error: %v", err)
		return err
	}

	return nil
}

// PostToBlog finds all GitHub issues within the latest numOfDay days and returns their titles and comments.
func (b *BookmarkMgr) PostToBlog(numOfDay string) ([]*github.Issue, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: b.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Parse the numOfDay string to an integer
	days, err := time.ParseDuration(numOfDay + "d")
	if err != nil {
		log.Printf("Error parsing number of days: %v", err)
		return nil, err
	}

	// Calculate the time for the start of the search window
	since := time.Now().Add(-days)

	// List all issues for the repository since the calculated time
	opts := &github.IssueListByRepoOptions{
		Since:       since,
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// Retrieve issues
	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, b.User, b.Repo, opts)
		if err != nil {
			log.Printf("Issues.ListByRepo returned error: %v", err)
			return nil, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Create a slice to hold the titles and comments of the issues
	var issueDetails []*github.Issue

	// Loop through the issues and add their titles and comments to the slice
	for _, issue := range allIssues {
		issueDetails = append(issueDetails, issue)
	}

	return issueDetails, nil
}
