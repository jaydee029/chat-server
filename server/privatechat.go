package main

type chatRooms struct {
	Client     map[*Client]bool
	register   chan *Client
	unregister chan *Client
	Broadcast  chan []byte
}

func (ws *Wserver) NewchatRooms() *chatRooms {
	return &chatRooms{
		Client:     make(map[*Client]bool, 2),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (room *chatRooms) RunchatRooms() {
	for {
		select {
		case client := <-room.register:
			room.registerClientinRoom(client)

		case client := <-room.unregister:
			room.unregisterClientinRoom(client)

		case msg := <-room.Broadcast:
			room.BroadcasttoRoom(msg)
		}
	}
}

func (room *chatRooms) registerClientinRoom(client *Client) {
	room.Client[client] = true
}

func (room *chatRooms) unregisterClientinRoom(client *Client) {
	delete(room.Client, client)
}

func (room *chatRooms) BroadcasttoRoom(msg []byte) {
	for client := range room.Client {
		client.sendto <- msg
	}
}
