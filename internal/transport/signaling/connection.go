// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package signaling provides utilities for WebRTC signaling
package signaling

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// WriteWait is the time allowed to write a message to the peer.
	WriteWait = 10 * time.Second
	// PongWait is the time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second
	// PingPeriod sends pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10
	// MaxMessageSize is the maximum message size allowed from peer.
	MaxMessageSize int64 = 64 * 1024
)

type Connection struct {
	ws     *websocket.Conn
	send   chan []byte
	hub    *Hub
	closed bool
}

// Send sends a message to the underlying connection
func (c *Connection) Send(m *Message) {
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	c.send <- b
}

// Close closes the connection
func (c *Connection) Close() {
	if !c.closed {
		if err := c.ws.Close(); err != nil {
			fmt.Println(err)
			log.Println("INFO [Router] Websocket was already closed:", err)
		}
		close(c.send)
		c.closed = true
	}
}

// ListenRead listens for incoming messages
func (c *Connection) ListenRead() {
	defer func() {
		c.hub.unregister <- c
	}()
	c.ws.SetReadLimit(MaxMessageSize)
	if err := c.ws.SetReadDeadline(time.Now().Add(PongWait)); err != nil {
		log.Println("WARNING [Router] Failed to set socket read deadline:", err)
	}
	c.ws.SetPongHandler(func(string) error {
		return c.ws.SetReadDeadline(time.Now().Add(PongWait))
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("ERROR [Router] Read message error:", err)
			break
		}

		m := &Message{connection: c}
		if err := json.Unmarshal(message, m); err != nil {
			log.Println("ERROR [Router] Invalid data sent for message:", string(message))
			continue
		}
		c.hub.Broadcast <- m
	}
}

// ListenWrite sends messages send to this connection
func (c *Connection) ListenWrite() {
	write := func(mt int, payload []byte) error {
		if err := c.ws.SetWriteDeadline(time.Now().Add(WriteWait)); err != nil {
			return err
		}
		return c.ws.WriteMessage(mt, payload)
	}
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if err := write(websocket.CloseMessage, []byte{}); err != nil {
					log.Println("WARNING [Router] Socket already closed:", err)
				}
				return
			}
			if err := write(websocket.TextMessage, message); err != nil {
				log.Println("WARNING [Router] Failed to write socket message:", err)
				return
			}
		case <-ticker.C:
			if err := write(websocket.PingMessage, nil); err != nil {
				log.Println("WARNING [Router] Failed to ping socket:", err)
				return
			}
		}
	}
}
