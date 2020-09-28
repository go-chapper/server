// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package broadcast provides utilities to broadcast messages
// Inspired by https://github.com/DATA-DOG/golang-websocket-hub and
// https://github.com/gorilla/websocket/tree/master/examples/chat
package broadcast

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ReadBufferSize  int = 1024
	WriteBufferSize int = 1024
)

// Hub describes an bradcasting hub interface
type Hub interface {
	Run()
	Register(*Connection)
	Unregister(*Connection)
	Broadcast(*Message)
	Token(string) (string, error)
	CreateConnection(http.ResponseWriter, *http.Request) (*Connection, error)
}

// SignalingHub is a broadcasting hub to deliver WebRTC signaling messages
type SignalingHub struct {
	sync.Mutex

	tokens map[string]string      // Auth token lookup
	conns  map[*Connection]string // Active connections
	peers  map[string]*Peer       // Active peers

	wsFactory  websocket.Upgrader // Websocket factory
	register   chan *Connection   // Register channel
	unregister chan *Connection   // Unregister channel

	broadcast chan *Message // Broadcast channel to distribute messages
}

// MessagingHub is a broadcasting hub to deliver real time chat messages
type MessagingHub struct {
	sync.Mutex

	tokens map[string]string      // Auth token lookup
	conns  map[*Connection]string // Active connections
	peers  map[string]*Peer       // Active peers

	wsFactory  websocket.Upgrader // Websocket factory
	register   chan *Connection   // Register channel
	unregister chan *Connection   // Unregister channel

	broadcast chan *Message // Broadcast channel to distribute messages
}

// Peer descripes one peer websocket connection
type Peer struct {
	Username   string
	Token      string
	connection *Connection
}

// NewSignalingHub returns a new signaling hub
func NewSignalingHub() *SignalingHub {
	h := &SignalingHub{
		tokens:     make(map[string]string),
		conns:      make(map[*Connection]string),
		peers:      make(map[string]*Peer),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		broadcast:  make(chan *Message),
	}

	h.wsFactory = websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}

	return h
}

// NewMessagingHub returns a new messaging hub
func NewMessagingHub() *MessagingHub {
	h := &MessagingHub{
		tokens:     make(map[string]string),
		conns:      make(map[*Connection]string),
		peers:      make(map[string]*Peer),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		broadcast:  make(chan *Message),
	}

	h.wsFactory = websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}

	return h
}
