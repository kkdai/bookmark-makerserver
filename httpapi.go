// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// IncomingMsg :
type IncomingMsg struct {
	User string `json:"User"`
	Repo string `json:"Repo"`
	Msg  string `json:"Msg"`
}

func bookmarkPost(w http.ResponseWriter, req *http.Request) {
	var in IncomingMsg

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Data read error:", err)
		return
	}
	// Replace newline characters with escaped version
	body = bytes.ReplaceAll(body, []byte("\n"), []byte("\\n"))

	log.Println("Body:", string(body))
	err = json.Unmarshal(body, &in)
	log.Println("Get request:", in)
	if err != nil {
		for name, values := range req.Header {
			// Loop over all values for the name.
			for _, value := range values {
				log.Println(name, value)
			}
		}

		log.Println("json unmarkshal error:", err, " body:", string(body))
		return
	}

	//Pass parameter
	githubToken := os.Getenv("GITHUB_TOKEN")

	bm := NewBookmark(in.User, in.Repo, githubToken)
	err = bm.SaveBookmark(in.Msg)
	if err != nil {
		log.Println("err=", err)
	}
	log.Println("Github issue post success! Msg=", in.Msg)
}

func serveHttpAPI(port string, existC chan bool) {
	go func() {
		if err, ok := <-existC; ok {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", bookmarkPost)
	http.ListenAndServe(":"+port, mux)
}
