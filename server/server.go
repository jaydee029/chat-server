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

type chatRooms struct {
	Client map[*Client]bool
}

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
	name := r.URL.Query().Get("name")
	client := newClient(name, conn, ws)

	go client.ReadInput()
	go client.WriteInput()

	fmt.Println("New client joined the chat server", client)

	ws.Register <- client

}

func (ws *Wserver) registerClient(client *Client) {
	room, ok := ws.ChatRooms[client.roomid]
	if ok {
		room.Client[client] = true
		msg := &Message{
			content: []byte(client.username + "has joined the chat"),
			roomid:  client.roomid,
			sender:  client.username,
		}
		client.Message <- msg
	}
	if !ok {
		ws.ChatRooms[client.roomid] = &chatRooms{
			Client: make(map[*Client]bool),
		}
		room := ws.ChatRooms[client.roomid]
		room.Client[client] = true
	}
}

func (ws *Wserver) unregisterClient(client *Client) {
	if len(ws.ChatRooms[client.roomid].Client) == 0 {
		delete(ws.ChatRooms[client.roomid].Client, client)
		delete(ws.ChatRooms, client.roomid)
	}

	msg := &Message{
		content: []byte(client.username + "has left the chat"),
		roomid:  client.roomid,
		sender:  client.username,
	}
	client.Message <- msg
	delete(ws.ChatRooms[client.roomid].Client, client)

}

func (ws *Wserver) BroadcastMessage(msg *Message) {
	for client := range ws.ChatRooms[msg.roomid].Client {
		if client.username != msg.sender {
			client.Message <- msg
		}
	}
}
