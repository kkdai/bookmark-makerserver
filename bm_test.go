package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gernest/mention"
	"github.com/mvdan/xurls"
)

func TestTweetHashTag(t *testing.T) {
	message := "how does it feel to be rejected? # it is #loner tt ggg sjdsj dj  #linker"
	tags := mention.GetTags('#', strings.NewReader(message))
	fmt.Println(tags)

	message2 := "hello @gernest I would like to @ follow you on twitter"
	tags2 := mention.GetTags('@', strings.NewReader(message2))
	fmt.Println(tags2)

	message3 := "hello @gernest I would  * like to #follow #go #ilike "
	tags3 := mention.GetTags('#', strings.NewReader(message3))
	content := mention.GetTags('*', strings.NewReader(message3))
	fmt.Println("tag:", tags3)
	fmt.Println("body:", content)
}

func TestStringWebLink(t *testing.T) {
	str1 := xurls.Relaxed.FindString("Do gophers live in golang.org?")
	if str1 == "" {
		t.Errorf("Cannot find string")
	}
	fmt.Println(str1)
	// "golang.org"
	str2 := xurls.Relaxed.FindAllString("foo.com is https://foo.com/.", -1)
	if str2 == nil {
		t.Errorf("Cannot find string")
	}
	fmt.Println(str2)
	// ["foo.com", "http://foo.com/"]
	str3 := xurls.Strict.FindAllString("foo.com is https://foo.com/.", -1)
	if str3 == nil {
		t.Errorf("Cannot find string")
	}
	fmt.Println(str3)
}

// TestIssueList :
func TestIssueList(t *testing.T) {
	// You need get your github token from https://github.com/settings/tokens

	token := os.Getenv("Token")
	user := os.Getenv("User")
	repo := os.Getenv("Repo")

	if len(token) == 0 {
		t.Skip("no token")
		return
	}

	t.Log(token, user, repo)
	testString := "Stateless datacenter load-balancing with Beamer | the morning paper https://t.co/0GFghfriwB"

	bm := NewBookmark(user, repo, token)
	err := bm.SaveBookmark(testString)
	if err != nil {
		t.Error(err)
	}
}

// TestPostToBlog tests the PostToBlog function with a real GitHub API call.
func TestPostToBlog(t *testing.T) {
	// Set up the BookmarkMgr with real token, user, and repo values.
	// WARNING: Do not hardcode tokens in your code; this is for demonstration purposes only.
	// Use environment variables or a secure method to handle tokens.
	token := os.Getenv("GITHUB_TOKEN")
	user := os.Getenv("User")
	repo := os.Getenv("Repo")

	if token == "" {
		t.Skip("Skipping test because GITHUB_TOKEN is not set")
	}

	bm := NewBookmark(user, repo, token)

	// Call the PostToBlog function with a number of days.
	issues, err := bm.PostToBlog("7") // Check the last 7 days.
	if err != nil {
		t.Errorf("PostToBlog returned an error: %v", err)
	}

	// Check if the issues slice is not nil.
	if issues == nil {
		t.Error("Expected a non-nil slice of issues, got nil")
	}

	// Here you could add more assertions, such as checking if the issues are within the last 7 days.
	for _, issue := range issues {
		fmt.Println("Title: ", issue.Title, " issue time:", issue.CreatedAt)
		fmt.Printf("Body: %s\n\n", *issue.Body)
	}
}
