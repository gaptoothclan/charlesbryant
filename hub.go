package main

type Broadcast struct {
	client Client
	message *Message
}

type Hub struct {
	clients map[*Client]bool
	register chan *Client
	unregister chan *Client
	broadcast chan *Broadcast
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan *Broadcast),
	}
}

func (h *Hub) run(){
	for {
		select {
		case client := <- h.register:

			// Get positions of all clients
			for clientMsg := range h.clients {
				client.send <- clientMsg.message
			}	
			h.clients[client] = true

		case client := <- h.unregister:
			messageJson := client.message
			messageJson.Delete = true

			// Notify all clients
			for clientMsg := range h.clients {
				if client != clientMsg {
					clientMsg.send <- messageJson
				}
			}	

			if ok, _ := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)				
			}
		case broadcast := <- h.broadcast:
			for client := range h.clients {
				if *client != broadcast.client {
					client.send <- broadcast.message
				}
			}
		}
	}
}