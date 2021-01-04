// Copyright (c) 2021-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

import (
	"errors"
	"net/http"
	"sync"

	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/utils"

	"github.com/gorilla/websocket"
)

var (
	ReadBufferSize  int = 1024
	WriteBufferSize int = 1024
)

var (
	ErrInvalidToken             = errors.New("invalid-token")
	ErrInvalidMessageType       = errors.New("invalid-message-type")
	ErrMessageTypeAlreadyExists = errors.New("message-type-already-exists")
)

// Hub is a broadcasting hub to deliver real time chat messages
type Hub struct {
	sync.Mutex

	tokens   map[string]string         // Auth token lookup
	peers    map[string]*Peer          // Active peers
	messages map[string]func() Message // Map of registered messages

	wsFactory websocket.Upgrader // Websocket factory
	logger    *log.Logger        // Logger

	broadcast chan *Message // Broadcast channel to distribute messages
}

// NewHub returns a new messaging hub
func NewHub(logger *log.Logger) *Hub {
	h := &Hub{
		tokens:    make(map[string]string),
		peers:     make(map[string]*Peer),
		messages:  make(map[string]func() Message),
		broadcast: make(chan *Message),
		logger:    logger,
	}

	h.wsFactory = websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header["Origin"][0]
			if origin == "http://localhost:8080" || origin == "chapper://." {
				return true
			}
			return false
		},
	}

	return h
}

// Run registers all messages and starts the main loop
func (h *Hub) Run() error {
	err := h.RegisterMessages(AllMessages())
	if err != nil {
		return err
	}

	go func() {
		for {
			m := <-h.broadcast
			h.handleMessage(m)
		}
	}()
	return nil
}

// Broadcast broadcasts a message
func (h *Hub) Broadcast(m *Message) {
	h.broadcast <- m
}

// Send sends the message to one or more receivers
func (h *Hub) Send(m *Message) {
	// for _, receiver := range m.To {
	// 	to, ok := h.peers[receiver]
	// 	if !ok {
	// 		return
	// 	}

	// 	to.connection.Send(m)
	// }
}

// Token returns a cryptographically secure random string. If the generation fails, an
// error is returned
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

// AuthenticatePeer authenticates a peer undentified by username with the provided token.
// If the authentication fails, an error is returned
func (h *Hub) AuthenticatePeer(username, token string) error {
	if t, ok := h.tokens[username]; !ok || t != token {
		return ErrInvalidToken
	}
	return nil
}

// NewPeer creates, registers and returns a new peer. If opening the websocket connection
// fails, an error is returned
func (h *Hub) NewPeer(username string, w http.ResponseWriter, r *http.Request) (*Peer, error) {
	ws, err := h.wsFactory.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	peer := &Peer{
		Username: username,
		ws:       ws,
		hub:      h,
	}

	h.Lock()
	h.peers[username] = peer
	h.Unlock()

	return peer, nil
}

// RegisterMessages registers an array of messages
func (h *Hub) RegisterMessages(messages []Message) error {
	for _, message := range messages {
		_, exists := h.messages[message.Type()]
		if exists {
			return ErrMessageTypeAlreadyExists
		}
		h.messages[message.Type()] = message.New()
	}
	return nil
}
