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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	plurgo "github.com/kkdai/plurgo/plurkgo"
)

type IncomingMsg struct {
	ConsumerToken  string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	Msg            string
}

func plurkPost(w http.ResponseWriter, req *http.Request) {
	var in IncomingMsg

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Data read error:", err)
		return
	}
	log.Println("Body:", string(body))
	err = json.Unmarshal(body, &in)
	log.Println("Get request:", in)
	if err != nil {
		log.Println("json unmarkshal error:", err)
		return
	}

	//Pass parameter
	var plurkCred plurgo.PlurkCredentials
	plurkCred.AccessSecret = in.AccessSecret
	plurkCred.AccessToken = in.AccessToken
	plurkCred.ConsumerSecret = in.ConsumerSecret
	plurkCred.ConsumerToken = in.ConsumerToken

	//Access plurk token
	accessToken, _, err := plurgo.GetAccessToken(&plurkCred)
	var data = map[string]string{}
	data["content"] = in.Msg
	data["qualifier"] = "shares"
	result, err := plurgo.CallAPI(accessToken, "/APP/Timeline/plurkAdd", data)
	if err != nil {
		log.Println("failed: %v ret=%v", err, result)
		return
	}
	log.Println("Plurk post success! Msg=", in.Msg)
}

func serveHttpAPI(port string, existC chan bool) {
	go func() {
		if err, ok := <-existC; ok {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", plurkPost)
	http.ListenAndServe(":"+port, mux)
}
