package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	ws       *Wserver
	username string
	Message  chan *Message
	roomid   string
}

type Message struct {
	content []byte
	roomid  string
	sender  string
}

func newClient(name string, conn *websocket.Conn, wserver *Wserver) *Client {
	return &Client{
		conn:     conn,
		ws:       wserver,
		username: name,
		Message:  make(chan *Message, 10),
	}
}

func (client *Client) ReadInput() {
	defer func() {
		client.ws.Unregister <- client
		client.conn.Close()
	}()
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("unexpected close error: %v", err)
			break
		}

		message := &Message{
			roomid:  client.roomid,
			sender:  client.username,
			content: msg,
		}
		client.ws.Broadcast <- message
	}
}

func (client *Client) WriteInput() {
	defer func() {
		client.conn.Close()
	}()
	for {
		Msg := <-client.Message
		w, err := client.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Printf("Error while writing the message ; %v", err)
			break
		}

		w.Write(Msg.content)
		w.Write([]byte("\n"))

		err = w.Close()
		if err != nil {
			log.Println(w.Close().Error())
		}

		//case Msg := <-client.sendto:

	}
}
