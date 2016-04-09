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
	Content     content  `json:"content"`
	CreatedTime int64    `json:"createdTime"`
	EventType   string   `json:"eventType"`
	From        string   `json:"from"`
	FromChannel int64    `json:"fromChannel"`
	ID          string   `json:"id"`
	To          []string `json:"to"`
	ToChannel   int64    `json:"toChannel"`
}

type content struct {
	ContentMetadata map[string]string `json:"contentMetadata"`
	ContentType     int               `json:"contentType"`
	CreatedTime     int64             `json:"createdTime"`
	DeliveredTime   int64             `json:"deliveredTime"`
	From            string            `json:"from"`
	ID              string            `json:"id"`
	Location        string            `json:"location"`
	Seq             string            `json:"seq"`
	Text            string            `json:"text"`
	To              []string          `json:"to"`
	ToType          int               `json:"toType"`
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

	channelID := os.Getenv("LINE_BOT_CHANNEL_ID")
	log.Println("CHANNEL_ID:", channelID)

	goji.Get("/", index)
	goji.Post("/bot/callback", callback)
	goji.Serve()
}
