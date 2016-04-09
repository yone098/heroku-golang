package main

import (
	"bytes"
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

type requestContent struct {
	To        []string `json:"to"`
	ToChannel int64    `json:"toChannel"`
	EventType string   `json:"eventType"`
	Content   content  `json:"content"`
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

	var reqContent requestContent
	for _, result := range m.Result {
		result.Content.Text = "ハゲ"
		reqContent = requestContent{
			To:        []string{result.Content.From},
			ToChannel: 1383378250,
			EventType: "138311608800106203",
			Content:   result.Content,
		}
	}
	log.Printf("### reqContent:%v\n", reqContent)

	b, err = json.Marshal(reqContent)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointURI := "https://trialbot-api.line.me/v1/events"
	req, err := http.NewRequest(
		"POST",
		endpointURI,
		bytes.NewBuffer(b),
	)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Line-ChannelID", channelID)
	req.Header.Set("X-Line-ChannelSecret", channelSecret)
	req.Header.Set("X-Line-Trusted-User-With-ACL", channelMID)

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("###STATUS:", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("body:", string(body))

	fmt.Fprintf(w, "OK")
}

var channelID string
var channelSecret string
var channelMID string

// main function
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	flag.Set("bind", ":"+port)

	channelID = os.Getenv("LINE_BOT_CHANNEL_ID")
	log.Println("CHANNEL_ID:", channelID)
	channelSecret = os.Getenv("LINE_BOT_CHANNEL_SECRET")
	log.Println("CHANNEL_SECRET:", channelSecret)
	channelMID = os.Getenv("LINE_BOT_CHANNEL_MID")
	log.Println("CHANNEL_MID:", channelMID)

	goji.Get("/", index)
	goji.Post("/bot/callback", callback)
	goji.Serve()
}
