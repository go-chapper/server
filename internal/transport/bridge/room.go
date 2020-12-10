// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bridge

import "encoding/json"

// Room maintains a number of active connections and handles messages from and to clients
type Room struct {
	name      string
	users     map[string]*User
	broadcast chan Message
	join      chan *User
	leave     chan *User
}

// Message carries a message with the additional info of the sender. A message will be
// delivered to everyone except the sender
type Message struct {
	data []byte
	user *User
}

// NewRoom creates and returns a new room
func NewRoom(roomHash string) *Room {
	return &Room{
		name:      roomHash,
		users:     make(map[string]*User),
		broadcast: make(chan Message),
		leave:     make(chan *User),
		join:      make(chan *User),
	}
}

// Run runs the main event loop
func (r *Room) Run() {
	go func() {
		for {
			select {
			case user := <-r.join:
				r.users[user.ID] = user
				go r.BroadcastEventJoin(user)
			case user := <-r.leave:
				_, exists := r.users[user.ID]
				if !exists {
					continue
				}

				delete(r.users, user.ID)
				close(user.send)
				r.BroadcastEventLeave(user)
			case message := <-r.broadcast:
				for _, user := range r.users {
					if message.user != nil && message.user.ID == user.ID {
						continue
					}
					user.send <- message.data
				}
			}
		}
	}()
}

// GetUsersList returns the user map as a slice of users
func (r *Room) GetUsersList() []*User {
	users := []*User{}

	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}

// GetParticipants returns all users except the provided user as a slice
func (r *Room) GetParticipants(u *User) []*User {
	users := []*User{}

	for _, user := range r.users {
		if u.ID == user.ID {
			continue
		}
		users = append(users, user)
	}

	return users
}

// Join connects the user to the room
func (r *Room) Join(user *User) {
	r.join <- user
}

// Leave disconnects the user from the room
func (r *Room) Leave(user *User) {
	r.leave <- user
}

// Broadcast sends a message to everyone except sender
func (r *Room) Broadcast(data []byte, user *User) {
	message := Message{data: data, user: user}
	r.broadcast <- message
}

// BroadcastEvent sends an event message to everyone except sender
func (r *Room) BroadcastEvent(event Event, user *User) error {
	json, err := json.Marshal(event)
	if err != nil {
		return err
	}

	r.Broadcast(json, user)
	return nil
}

// BroadcastEventJoin creates a join event and broadcasts it
func (r *Room) BroadcastEventJoin(user *User) error {
	event := Event{Type: EventJoin, User: user.ToPublic()}
	return r.BroadcastEvent(event, user)
}

// BroadcastEventLeave creates a leave event and broadcasts it
func (r *Room) BroadcastEventLeave(user *User) error {
	event := Event{Type: EventLeave, User: user.ToPublic()}
	return r.BroadcastEvent(event, user)
}

// BroadcastEventMute creates microphone mute event and broadcasts it
func (r *Room) BroadcastEventMute(user *User) error {
	event := Event{Type: EventMute, User: user.ToPublic()}
	return r.BroadcastEvent(event, user)
}

// BroadcastEventUnmute creates microphone unmute event and broadcasts it
func (r *Room) BroadcastEventUnmute(user *User) error {
	event := Event{Type: EventUnmute, User: user.ToPublic()}
	return r.BroadcastEvent(event, user)
}
