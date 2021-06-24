package chat

import (
	"encoding/json"
	"html"
	"net/http"
	"time"

	"github.com/lazy-void/chatapp/internal/models"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type contextKey string

var ContextUserKey = contextKey("user")

const (
	// Time allowed to write message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Sends pings to peer with this period. Must be less than pongWait.
	pingPeriod = (9 * pongWait) / 10

	// Maximum message size allowed from peer.
	messageMaxSize = 2048
)

// API actions.
const (
	loadMoreAction  = "loadMore"
	broadcastAction = "broadcast"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Request struct {
	Action string `json:"action"`

	// if client wants to broadcast a message
	Message string `json:"message"`

	// if client wants to load more messages
	Offset int `json:"offset"`
}

type Client struct {
	hub *Hub

	// Information about the user
	user models.User

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	sendMessage chan Message

	// Channel of outbound responses to requests
	sendResponse chan Response
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(messageMaxSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Err(err).Msg("error while setting read deadline on websocket connection")
		return
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Err(err).Msg("websocket connection closed unexpectedly")
			}
			return
		}

		var req Request
		err = json.Unmarshal(message, &req)
		if err != nil {
			log.Err(err).Msg("error unmarshalling client request")
			return
		}

		switch req.Action {
		case broadcastAction:
			c.hub.broadcast <- Message{
				Text:     html.EscapeString(req.Message),
				Username: c.user.Username,
				Created:  time.Now().UTC(),
			}
		case loadMoreAction:
			messages, err := c.hub.loadMore(req.Offset)
			if err != nil {
				log.Err(err).Msg("error loading messages from db")
				continue
			}

			c.sendResponse <- Response{Request: req, Messages: messages}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	writeMessage := func(msg []byte) error {
		err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err != nil {
			return err
		}

		err = c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return err
		}

		return nil
	}

	for {
		select {
		case message, ok := <-c.sendMessage:
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
			}

			// add all queued messages to the update
			update := Update{Messages: []Message{message}}
			for i := 0; i < len(c.sendMessage); i++ {
				update.Messages = append(update.Messages, <-c.sendMessage)
			}

			jsonUpdate, err := json.Marshal(update)
			if err != nil {
				log.Err(err).Msg("error marshaling update to json")
				return
			}

			err = writeMessage(jsonUpdate)
			if err != nil {
				log.Err(err).Msg("error sending message to websocket")
				return
			}
		case resp, ok := <-c.sendResponse:
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
			}

			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Err(err).Msg("error marshaling response to json")
				return
			}

			err = writeMessage(jsonResp)
			if err != nil {
				log.Err(err).Msg("error sending message to websocket")
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Err(err).Msg("error setting write deadline on websocket")
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Err(err).Msg("error writing message to websocket")
				return
			}
		}
	}
}

// ServeWS handles websocket requests from the peer.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("error upgrading connection to websocket")
		return
	}

	user := r.Context().Value(ContextUserKey).(models.User)
	c := &Client{
		hub:          hub,
		user:         user,
		conn:         conn,
		sendMessage:  make(chan Message, 256),
		sendResponse: make(chan Response),
	}
	c.hub.register <- c

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go c.readPump()
	go c.writePump()
}
