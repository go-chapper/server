// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package broadcast provides utilities to broadcast messages
package broadcast

import "fmt"

func (h *SignalingHub) handleRegister(c *Connection) {
	h.Lock()
	defer h.Unlock()

	h.conns[c] = true
}

func (h *SignalingHub) handleUnregister(c *Connection) {
	h.Lock()
	defer h.Unlock()

	_, ok := h.conns[c]
	if !ok {
		// h.log.Println("[WARN] cannot unregister connection, it is not registered.")
		return
	}

	// h.log.Println("[DEBUG] unregistering socket connection")
	c.Close()
	delete(h.conns, c)
	// TODO <2020/14/09>: Unregister peer, change h.conns to map[*Connection]string
	fmt.Println("Unregistered connection")
}

func (h *SignalingHub) handleMessage(m *Message) {
	switch m.Topic {
	case "subscribe":
		if m.Token == "" || m.Username == "" {
			return
		}

		t, ok := h.tokens[m.Username]
		if !ok || t != m.Token {
			go h.Unregister(m.connection)
			return
		}

		h.peers[m.Username] = &Peer{
			Username:   m.Username,
			Token:      m.Token,
			connection: m.connection,
		}
		fmt.Println("Peer added")
	case "new-chat-offer":
		// Offers callee to begin a new chat
		if m.To == "" || m.From == "" {
			return
		}

		to, ok := h.peers[m.To]
		if !ok {
			return
		}

		to.connection.Send(m)
	case "new-chat-answer":
		// Answers a new chat offer to caller
		if m.To == "" || m.From == "" {
			return
		}

		to, ok := h.peers[m.To]
		if !ok {
			return
		}

		to.connection.Send(m)
	case "text-offer":
		if m.InvalidSDP() {
			return
		}

		to, ok := h.peers[m.To]
		if !ok {
			return
		}

		to.connection.Send(m)
	case "new-ice-candidate":
		if m.InvalidSDP() {
			return
		}

		to, ok := h.peers[m.To]
		if !ok {
			return
		}

		to.connection.Send(m)
	}
}
