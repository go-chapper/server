// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package signaling provides utilities for WebRTC signaling
// Inspired by https://github.com/DATA-DOG/golang-websocket-hub and
// https://github.com/gorilla/websocket/tree/master/examples/chat
package signaling

import (
	"fmt"
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
	tokens      map[string]string           // Authentication tokens
	connections map[*Connection]*Subscriber // Active connections
	subscribers map[string]*Subscriber      // Active subscribers

	wsFactory  websocket.Upgrader // Websocket factory
	register   chan *Connection   // Register channel
	unregister chan *Connection   // Unregister channel
	subscribe  chan *Subscription // Subscribe channel

	Broadcast chan *Message // Broadcast channel to deliver messages
}

type Message struct {
	Data  interface{} `json:"data"`
	Topic string      `json:"topic"`
}

// New returns a new Hub
func New() *Hub {
	h := &Hub{
		tokens:      make(map[string]string),
		connections: make(map[*Connection]*Subscriber),
		subscribers: make(map[string]*Subscriber),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		subscribe:   make(chan *Subscription),
		Broadcast:   make(chan *Message),
	}

	factory := websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}
	h.wsFactory = factory

	return h
}

func (h *Hub) Run() {
	go func() {
		for {
			select {
			case c := <-h.register:
				fmt.Println(c)
				// h.doRegister(c)
			case c := <-h.unregister:
				fmt.Println(c)
				// h.doUnregister(c)
			case m := <-h.Broadcast:
				fmt.Println(m)
				// h.doBroadcast(m)
			case s := <-h.subscribe:
				fmt.Println(s.Username, s.Token, h.tokens[s.Username])
				// h.doSubscribe(s)
			}
		}
	}()
}

func (h *Hub) Register(c *Connection) {
	h.register <- c
}

func (h *Hub) Subscribe(s *Subscription) {
	h.subscribe <- s
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
