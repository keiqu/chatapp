package chat

import (
	"bytes"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Sends pings to peer with this period. Must be less than pongWait.
	pingPeriod = (9 * pongWait) / 10

	// Maximum message size allowed from peer.
	messageMaxSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		if err := c.conn.Close(); err != nil {
			log.Err(err)
		}
	}()

	c.conn.SetReadLimit(messageMaxSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Err(err)
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
				log.Err(err)
			}
			return
		}

		c.hub.broadcast <- bytes.TrimSpace(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			log.Err(err)
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Err(err)
				return
			}

			if !ok {
				// The hub closed the channel.
				err := c.conn.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Err(err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Err(err)
				return
			}

			if _, err := w.Write(message); err != nil {
				log.Err(err)
				return
			}

			if err := w.Close(); err != nil {
				log.Err(err)
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Err(err)
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Err(err)
				return
			}
		}
	}
}

// ServeWS handles websocket requests from the peer.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err)
		return
	}
	c := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	c.hub.register <- c

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go c.readPump()
	go c.writePump()
}
