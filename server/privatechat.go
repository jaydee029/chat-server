package main

type Chats struct {
	Client     map[*Client]bool
	register   chan *Client
	unregister chan *Client
	Broadcast  chan []byte
}

func (ws *Wserver) NewChat() {
	
}
