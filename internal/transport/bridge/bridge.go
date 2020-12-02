// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bridge provides utilities to create, keep track of and delete WebRTC calls
package bridge

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v2"
)

var (
	// ErrNoSuchRoom indicates no such room with given key exists in the pool of active
	// rooms
	ErrNoSuchRoom = errors.New("No such room")

	// ErrDuplicateRoom indicates the room with given key already exists
	ErrDuplicateRoom = errors.New("Duplicate room")
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connectionConfig = webrtc.Configuration{
	SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
}

// Bridge keeps track of active rooms
type Bridge struct {
	rooms map[string]*Room
}

// NewBridge returns a new bridge
func NewBridge() *Bridge {
	return &Bridge{
		rooms: make(map[string]*Room),
	}
}

// GetRoom returns room with 'roomHash' from the pool of active rooms. If the room doesn't
// exists ErrNoSuchRoom is returned
func (b *Bridge) GetRoom(roomHash string) (*Room, error) {
	room, exists := b.rooms[roomHash]
	if !exists {
		return nil, ErrNoSuchRoom
	}

	return room, nil
}

// GetOrCreateRoom returns an existing room, or if no such room exists, creates a new one
// and returns it
func (b *Bridge) GetOrCreateRoom(roomHash string) (*Room, error) {
	room, exists := b.rooms[roomHash]
	if exists {
		return room, nil
	}

	return b.AddRoom(roomHash)
}

// AddRoom adds a room to the pool of active rooms with the key 'roomHash'. If the room
// already exists ErrDuplicateRoom is returned
func (b *Bridge) AddRoom(roomHash string) (*Room, error) {
	_, exists := b.rooms[roomHash]
	if exists {
		return nil, ErrDuplicateRoom
	}

	room := NewRoom(roomHash)

	go room.Run()
	b.rooms[roomHash] = room
	return room, nil
}

// RemoveRoom removes a room with key 'roomHash' from the pool of active rooms. If the
// room with the given key doesn't exist ErrNoSuchRoom is returned
func (b *Bridge) RemoveRoom(roomHash string) error {
	_, exists := b.rooms[roomHash]
	if !exists {
		return ErrNoSuchRoom
	}

	delete(b.rooms, roomHash)
	return nil
}

func (b *Bridge) Connect(username, roomHash string, w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Setup WebRTC session
	mediaEngine := webrtc.MediaEngine{}
	mediaEngine.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
	pc, err := api.NewPeerConnection(connectionConfig)
	if err != nil {
		return
	}

	room, err := b.GetOrCreateRoom(roomHash)
	if err != nil {
		return
	}

	user := NewUser(username, conn, pc, room)
	user.AddListeners()

	go user.startRead()
	go user.startWrite()
}
