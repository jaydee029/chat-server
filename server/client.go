package main

import (
	"log"	
	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	ws     *Wserver
	sendto chan []byte
}

func newClient(conn *websocket.Conn, wserver *Wserver) *Client {
	return &Client{
		conn:   conn,
		ws:     wserver,
		sendto: make(chan []byte),
	}
}

func (client *Client) ReadInput() {
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("unexpected close error: %v", err)
			break
		}

		client.ws.Broadcast <- msg
	}
}

func (client *Client) WriteInput() {
	for {
		Msg := <-client.sendto
		w, err := client.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Printf("Error while writing the message ; %v", err)
			break
		}

		w.Write(Msg)

		n := len(client.sendto)
		for i := 0; i < n; i++ {
			w.Write([]byte("\n"))
			w.Write(<-client.sendto)

		}

		err = w.Close()
		if err != nil {
			log.Println(w.Close().Error())
		}

		//case Msg := <-client.sendto:

	}
}
