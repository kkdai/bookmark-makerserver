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

//TestIssueList :
func TestIssueList(t *testing.T) {
	// You need get your github token from https://github.com/settings/tokens

	token := os.Getenv("Token")
	user := os.Getenv("User")
	repo := os.Getenv("Repo")

	testString := "mvdan/xurls: Extract urls from text help you to find link in string in #golang https://github.com/mvdan/xurls"

	bm := NewBookmark(user, repo, token)
	err := bm.SaveBookmark(testString)
	if err != nil {
		t.Error(err)
	}
}
