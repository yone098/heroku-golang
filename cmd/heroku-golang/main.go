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

func hoge(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", "hoge")
}

func recieve(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
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

func linebot(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is Test, %s!", "LINE BOT")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var m msg
	if err := json.Unmarshal(b, &m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("msg:%v\n", m)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	flag.Set("bind", ":"+port)

	goji.Get("/test", hoge)
	goji.Post("/hello/:name", recieve)
	goji.Post("/linebot", linebot)
	goji.Serve()
}
