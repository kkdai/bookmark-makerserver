package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

type ResponseJson struct {
	Summary     string   `json:"summary"`
	FullContent string   `json:"full_content"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
}

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

	//Promot
	prompt := "根據這份資料，將它擷取出有用資料如下 summary, full_content, url, tags 給我 json ，其中 url 要從 msg 裡面擷取出來,  tags 除了資料中外，根據內容也整理出相關 tags ----"

	// Using Gemini
	ret := GeminiChatComplete(prompt + tweet)
	log.Println("ret:", ret)

	// Remove first and last line,	which are the backticks.
	jsonData := removeFirstAndLastLine(ret)
	log.Println("jsonData:", jsonData)

	// Parse json and insert NotionDB
	var res ResponseJson
	err := json.Unmarshal([]byte(jsonData), &res)
	if err != nil {
		log.Println("Error parsing JSON:", err)
	}
	log.Println(res)

	// Create a GitHub issue.
	input := &github.IssueRequest{
		Title:  &res.Summary,
		Body:   github.String(jsonData),
		Labels: &res.Tags,
	}

	_, _, err = client.Issues.Create(ctx, b.User, b.Repo, input)
	if err != nil {
		log.Printf("Issues.Create returned error: %v", err)
		return err
	}

	return nil
}

// PostToBlog finds all GitHub issues since the provided time without the "archived" label.
// If modifyLabels is true, it adds the "archived" label to them.
// It returns their titles and comments.
func (b *BookmarkMgr) PostToBlog(since time.Time, modifyLabels bool) ([]*github.Issue, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: b.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Define the label to be excluded and potentially added
	labelToExclude := "archived"

	// List all issues for the repository since the provided time, not closed, and without the "archived" label.
	opts := &github.IssueListByRepoOptions{
		Since: since,
		State: "open",
		// Labels:      []string{"-" + labelToExclude},
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

	// Loop through the issues and add the "archived" label if modifyLabels is true.
	if modifyLabels {
		for _, issue := range allIssues {
			// Add the "archived" label to the issue
			_, _, err := client.Issues.AddLabelsToIssue(ctx, b.User, b.Repo, *issue.Number, []string{labelToExclude})
			if err != nil {
				log.Printf("Failed to add 'archived' label to issue #%d: %v", *issue.Number, err)
				// Decide if you want to return an error or just log it
				// return nil, err
			} else {
				log.Printf("Added 'archived' label to issue #%d successfully", *issue.Number)
			}
		}
	}

	// Return all retrieved issues, regardless of whether labels were modified.
	return allIssues, nil
}
