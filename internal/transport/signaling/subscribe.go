// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package signaling provides utilities for WebRTC signaling
package signaling

type Subscriber struct {
	connections map[*Connection]bool
	Username    string
}

type Subscription struct {
	Username   string
	Token      string
	connection *Connection
}
