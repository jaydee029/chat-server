package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *client {
	return &client{
		conn: conn,
	}

}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat", nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from server:", err)
				return
			}

			writer := bufio.NewWriter(os.Stdout)
			writer.Write(msg)
			writer.Flush()
			//fmt.Println("Message from server:", string(msg))
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-interrupt:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Print("Enter message: ")
			msg, _ := reader.ReadString('\n')
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Error sending message:", err)
				return
			}
			time.Sleep(100 * time.Millisecond) // To prevent flooding the server
		}
	}
}
