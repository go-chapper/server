// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package signaling provides utilities for WebRTC signaling
// Inspired by https://github.com/DATA-DOG/golang-websocket-hub and
// https://github.com/gorilla/websocket/tree/master/examples/chat
package signaling

import (
	"net/http"
	"sync"

	"chapper.dev/server/internal/utils"

	"github.com/gorilla/websocket"
)

var (
	ReadBufferSize  int = 1024
	WriteBufferSize int = 1024
)

type Hub struct {
	sync.Mutex

	tokens map[string]string    // Auth token lookup
	conns  map[*Connection]bool // Active connections
	peers  map[string]*Peer     // Active peers

	wsFactory  websocket.Upgrader // Websocket factory
	register   chan *Connection   // Register channel
	unregister chan *Connection   // Unregister channel

	Broadcast chan *Message // Broadcast channel to distribute messages
}

type Peer struct {
	Username   string
	Token      string
	connection *Connection
}

func New() *Hub {
	h := &Hub{
		tokens:     make(map[string]string),
		conns:      make(map[*Connection]bool),
		peers:      make(map[string]*Peer),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		Broadcast:  make(chan *Message),
	}

	h.wsFactory = websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}

	return h
}

func (h *Hub) Run() {
	go func() {
		for {
			select {
			case c := <-h.register:
				h.handleRegister(c)
			case c := <-h.unregister:
				h.handleUnregister(c)
			case m := <-h.Broadcast:
				h.handleMessage(m)
			}
		}
	}()
}

// Register registers a new connection
func (h *Hub) Register(c *Connection) {
	h.register <- c
}

// Register registers a new connection
func (h *Hub) Unregister(c *Connection) {
	h.unregister <- c
}

// Token returns a cryptographically secure random string
func (h *Hub) Token(key string) (string, error) {
	h.Lock()
	defer h.Unlock()

	if t, ok := h.tokens[key]; ok {
		return t, nil
	}

	s, err := utils.RandomCryptoString(16)
	if err != nil {
		return "", err
	}

	h.tokens[key] = s
	return s, nil
}

// CreateConnection creates a new connection and returns it
func (h *Hub) CreateConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	ws, err := h.wsFactory.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &Connection{
		send: make(chan []byte, 256),
		ws:   ws,
		hub:  h,
	}, nil
}
