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

//TestIssueList :
func TestIssueList(t *testing.T) {
	// You need get your github token from https://github.com/settings/tokens

	token := os.Getenv("Token")
	user := os.Getenv("User")
	repo := os.Getenv("Repo")

	testString := "Facebook开源AI对话研究平台ParlAI ，解决人机对话最常见5类问题 # #curate #feedly  感覺不錯玩，可以玩玩看。 http://36kr.com/p/5075718.html?ktm_source=feed"

	bm := NewBookmark(user, repo, token)
	err := bm.SaveBookmark(testString)
	if err != nil {
		t.Error(err)
	}
}
