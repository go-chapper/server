// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package signaling provides utilities to handle WebRTC signaling
package signaling

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"chapper.dev/server/internal/modules/jwt"

	"github.com/gorilla/websocket"
)

var (
	ReadBufferSize  int = 1024
	WriteBufferSize int = 1024
)

// Hub is the central message hub. It is responsible for receiving, forwarding and sending
// messages from and to connected clients.
type Hub struct {
	sync.Mutex

	tokens    map[string]string  // Auth token lookup
	clients   map[string]*Client // Active client connections
	wsFactory websocket.Upgrader // Websocket factory
	broadcast chan Message       // Incoming client messages

	messages map[string]Message // Map of registered messages
}

// TypedMessage is an in- or outgoing message
type TypedMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewHub creates a new hub and returns it
func NewHub() *Hub {
	return &Hub{
		tokens:  make(map[string]string),
		clients: make(map[string]*Client),
		wsFactory: websocket.Upgrader{
			ReadBufferSize:  ReadBufferSize,
			WriteBufferSize: WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("origin")
				if origin == "http://localhost:8080" || origin == "chapper://." {
					return true
				}
				return false
			},
		},
		broadcast: make(chan Message),
	}
}

// Run runs the hub. It listens for incoming messages in it's own go routine. If there was
// en error while handling the message the client connection gets automatically closed.
func (h *Hub) Run() {
	go func() {
		for {
			msg := <-h.broadcast
			if err := msg.Handle(h); err != nil {
				// log error
			}
		}
	}()
}

func (h *Hub) RegisterMessages(messages map[string]Message) {
	h.messages = messages
}

// Upgrade upgrades the connection to a websocket connection and starts the read and write
// listeners
func (h *Hub) Upgrade(w http.ResponseWriter, req *http.Request, claims *jwt.Claims) error {
	conn, err := h.wsFactory.Upgrade(w, req, nil)
	if err != nil {
		return err
	}

	client := NewClient(conn, h, claims.Username)

	go client.listenRead(20 * time.Second)
	go client.listenWrite(5 * time.Second)
	return nil
}

func (h *Hub) GetClient(key string) *Client {
	return h.clients[key]
}

// readIncoming reads an incoming typed message which will get marshaled into the final
// event
func (h *Hub) readIncoming(r io.Reader) (Message, error) {
	incoming := TypedMessage{}
	if err := json.NewDecoder(r).Decode(&incoming); err != nil {
		return nil, err
	}

	event, ok := h.messages[incoming.Type]
	if !ok {
		return nil, fmt.Errorf("failed to read incoming message: unknow type %s", incoming.Type)
	}

	if err := json.Unmarshal(incoming.Payload, &event); err != nil {
		return nil, fmt.Errorf("failed to read incoming message: failed to unmarshal: %v", err)
	}
	return event, nil
}

// writeOutgoing writes an outgoing typed message based on a message
func (h *Hub) writeOutgoing(outgoing Message) (TypedMessage, error) {
	payload, err := json.Marshal(outgoing)
	if err != nil {
		return TypedMessage{}, err
	}
	return TypedMessage{
		Type:    outgoing.Type(),
		Payload: payload,
	}, nil
}
