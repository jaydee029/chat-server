package main

import (
	"log"
	"net/http"
)

type Wserver struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	ChatRooms  map[*chatRooms]bool
}

func main() {

	Wsserver := &Wserver{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		ChatRooms:  make(map[*chatRooms]bool),
	}

	http.HandleFunc("/chat", Wsserver.handleChat)
	log.Printf("The server is live on port %s\n", ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
