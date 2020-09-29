// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package broadcast provides utilities to broadcast messages
package broadcast

func (h *SignalingHub) handleRegister(c *Connection) {
	h.Lock()
	defer h.Unlock()

	h.conns[c] = ""
}

func (h *SignalingHub) handleUnregister(c *Connection) {
	h.Lock()
	defer h.Unlock()

	username, ok := h.conns[c]
	if !ok {
		return
	}

	c.Close()
	delete(h.conns, c)
	delete(h.peers, username)
}

func (h *SignalingHub) handleMessage(m *Message) {
	switch m.Type {
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
		h.conns[m.connection] = m.Username
	case "new-chat-offer":
		// Offers callee to begin a new chat
		if len(m.To) == 0 || m.From == "" {
			return
		}

		h.Send(m)
	case "new-chat-answer":
		// Answers a new chat offer to caller
		if len(m.To) == 0 || m.From == "" {
			return
		}

		h.Send(m)
	case "text-offer":
		if m.InvalidSDP() {
			return
		}

		h.Send(m)
	case "new-ice-candidate":
		if m.InvalidSDP() {
			return
		}

		h.Send(m)
	}
}
