package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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
	sinceTime := time.Now().Add(-3 * 24 * time.Hour) // Example: 3 days ago
	shouldModifyLabels := false                      // Set to false if you do not want to modify labels
	issues, err := bm.PostToBlog(sinceTime, shouldModifyLabels)
	if err != nil {
		// handle error
		t.Error(err)
	}
	// process issues
	// Here you could add more assertions, such as checking if the issues are within the last 7 days.
	for _, issue := range issues {
		t.Log("Title: ", *issue.Title, " issue time:", issue.CreatedAt)
		t.Log("Body:", *issue.Body)
		t.Log("issue link:", issue.GetHTMLURL())
	}
}

// TestScrapeURL tests the scrapeURL function with a real API call.
func TestScrapeURL(t *testing.T) {
	// Load environment variables.
	fc_token := os.Getenv("FC_AUTH_TOKEN")
	if fc_token == "" {
		t.Skip("Skipping test because FC_AUTH_TOKEN is not set")
	}

	// Set up the URL to scrape.
	url := "https://developers.googleblog.com/en/gemini-15-flash-updates-google-ai-studio-gemini-api/"

	// Call the scrapeURL function.
	resp, err := scrapeURL(url)
	if err != nil {
		t.Error(err)
	}

	// Print the response.
	t.Log("Title:", resp.Data.Metadata.Title)
	t.Log("Description:", resp.Data.Metadata.Description)
	t.Log("Source URL:", resp.Data.Metadata.SourceURL)
	t.Log("Content:", resp.Data.Content)
}