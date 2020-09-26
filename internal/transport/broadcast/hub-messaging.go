// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package broadcast provides utilities to broadcast messages
package broadcast

import (
	"net/http"

	"chapper.dev/server/internal/utils"
)

func (h *MessagingHub) Run() {
	go func() {
		for {
			select {
			case c := <-h.register:
				h.handleRegister(c)
			case c := <-h.unregister:
				h.handleUnregister(c)
			case m := <-h.broadcast:
				h.handleMessage(m)
			}
		}
	}()
}

// Register registers a new connection
func (h *MessagingHub) Register(c *Connection) {
	h.register <- c
}

// Register registers a new connection
func (h *MessagingHub) Unregister(c *Connection) {
	h.unregister <- c
}

// Broadcast broadcasts a message
func (h *MessagingHub) Broadcast(m *Message) {
	h.broadcast <- m
}

// Token returns a cryptographically secure random string
func (h *MessagingHub) Token(key string) (string, error) {
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
func (h *MessagingHub) CreateConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
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
