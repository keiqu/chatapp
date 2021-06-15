// Package chat implements logic of a websocket chat.
package chat

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	//  Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from the clients.
	unregister chan *Client
}

// NewHub initializes new instance of the Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts chat Hub. Messages are sent to out channel for
// them to be stored in persistent storage.
func (h *Hub) Run(out chan<- string) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if h.clients[client] {
				close(client.send)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			out <- string(message)

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
