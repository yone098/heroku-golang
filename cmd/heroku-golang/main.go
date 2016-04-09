package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func index(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Println("called hoge")
	fmt.Errorf("This is Test, %s!", "LINE BOT")
	fmt.Fprintf(w, "Hello %s!", "hoge")
}

type callbackMessage struct {
	Result []msg `json:"result"`
}

type msg struct {
	ID              string                 `json:"id"`
	ContentType     int                    `json:"conetntType"`
	From            string                 `json:"from"`
	CreatedTime     int                    `json:"createdTime"`
	To              []string               `json:"to"`
	ToType          int                    `json:"toType"`
	ContentMetadata map[string]interface{} `json:"contentMetadata"`
	Text            string                 `json:"text"`
	Location        map[string]interface{} `json:"location"`
}

// callback function
func callback(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Println("called callback")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("msg:", string(b))

	var m callbackMessage
	if err := json.Unmarshal(b, &m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("msg:%v\n", m)
}

// main function
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	flag.Set("bind", ":"+port)

	goji.Get("/", index)
	goji.Post("/bot/callback", callback)
	goji.Serve()
}
