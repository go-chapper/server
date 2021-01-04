// Copyright (c) 2021-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

import "github.com/gorilla/websocket"

// Peer describes one client connected to the hub. Each peer has a unique username to
// identify itself, a token to authenticate and the underlying websocket connection for
// real-time communication
type Peer struct {
	Username string
	Token    string
	ws       *websocket.Conn
	hub      *Hub
}

// Authenticate authenticates a peer. If the authentication fails, an error is returned
func (p *Peer) Authenticate(token string) error {
	return p.hub.AuthenticatePeer(p.Username, token)
}
