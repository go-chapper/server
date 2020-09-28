// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package broadcast provides utilities to broadcast messages
// Inspired by https://github.com/DATA-DOG/golang-websocket-hub and
// https://github.com/gorilla/websocket/tree/master/examples/chat
package broadcast

// Message describes a message to be sent and received
type Message struct {
	Type       string   `json:"type,omitempty"`
	Username   string   `json:"username,omitempty"`
	Token      string   `json:"token,omitempty"`
	From       string   `json:"from,omitempty"`
	To         []string `json:"to,omitempty"`
	Data       string   `json:"data,omitempty"`
	SDP        string   `json:"sdp,omitempty"`
	Accepted   bool     `json:"accepted,omitempty"`
	connection *Connection
}

// InvalidSDP returns if the SDP is invalid
func (m *Message) InvalidSDP() bool {
	return len(m.To) == 0 || m.SDP == ""
}
