// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bridge provides utilities to create, keep track of and delete WebRTC calls
package bridge

import "errors"

var (
	// ErrNoSuchRoom indicates no such room with given key exists in the pool of active
	// rooms
	ErrNoSuchRoom = errors.New("No such room")

	// ErrDuplicateRoom indicates the room with given key already exists
	ErrDuplicateRoom = errors.New("Duplicate room")
)

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

func (b *Bridge) Upgrade() {

}
