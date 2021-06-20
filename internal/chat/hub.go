// Package chat implements logic of a websocket chat.
package chat

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/lazy-void/chatapp/internal/models"
)

type Message struct {
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
}

type Response struct {
	Request  Request   `json:"request"`
	Messages []Message `json:"messages"`
}

type Update struct {
	Messages []Message `json:"messages"`
}

type MessageInterface interface {
	Insert(string, time.Time) (int, error)
	Get(n, offset int) ([]models.Message, error)
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	//  Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from the clients.
	unregister chan *Client

	// Get and insert chat messages from/in the storage.
	messages MessageInterface
}

// NewHub initializes new instance of the Hub.
func NewHub(messages MessageInterface) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		messages:   messages,
	}
}

// Run starts chat Hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if h.clients[client] {
				close(client.sendMessage)
				close(client.sendResponse)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			err := h.save(message)
			if err != nil {
				log.Err(err).Msg("error while saving message to db")
				continue
			}

			for client := range h.clients {
				select {
				case client.sendMessage <- message:
				default:
					close(client.sendMessage)
					close(client.sendResponse)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) save(m Message) error {
	_, err := h.messages.Insert(m.Text, m.Created)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hub) loadMore(offset int) ([]Message, error) {
	messages, err := h.messages.Get(100, offset)
	if err != nil {
		return []Message{}, err
	}
	if len(messages) == 0 {
		return []Message{}, nil
	}

	chatMessages := make([]Message, len(messages))
	for i, m := range messages {
		chatMessages[i] = Message{Text: m.Text, Created: m.Created}
	}

	return chatMessages, nil
}
