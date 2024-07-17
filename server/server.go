package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
)

func (ws *Wserver) Runserver() {

	for {
		select {
		case Client := <-ws.Register:
			ws.registerClient(Client)
		case Client := <-ws.Unregister:
			ws.unregisterClient(Client)
		case Message := <-ws.Broadcast:
			ws.BroadcastMessage(Message)
		}

	}
}

func (ws *Wserver) handleChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	go ws.Runserver()
	client := newClient(conn, ws)

	go client.ReadInput()
	go client.WriteInput()

	fmt.Println("New client joined the chat server", client)

	ws.Register <- client

}

func (ws *Wserver) registerClient(client *Client) {
	ws.Clients[client] = true
}

func (ws *Wserver) unregisterClient(client *Client) {
	delete(ws.Clients, client)
}

func (ws *Wserver) BroadcastMessage(msg []byte) {
	for client := range ws.Clients {
		client.sendto <- msg
	}
}
