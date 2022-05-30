package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Markdown map[string]string `json:"markdown"`
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	// r.GetBody()
	var message Message
	json.NewDecoder(r.Body).Decode(&message)
	log.Println(message)

}
func uploadMedia(w http.ResponseWriter, r *http.Request) {}
