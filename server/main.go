package main

import (
	"log"
	"net/http"
)

type Wserver struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	ChatRooms  map[string]*chatRooms
}

func main() {

	Wsserver := &Wserver{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		ChatRooms:  make(map[string]*chatRooms),
	}

	http.HandleFunc("/chat", Wsserver.handleChat)
	log.Printf("The server is live on port %s\n", ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
