// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package broadcast provides utilities to broadcast messages
package broadcast

func (h *MessagingHub) handleRegister(c *Connection) {
	h.Lock()
	defer h.Unlock()

	h.conns[c] = ""
}

func (h *MessagingHub) handleUnregister(c *Connection) {
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

func (h *MessagingHub) handleMessage(m *Message) {
	switch m.Type {
	case "status-writing":
		// Indicate a chat member is writing
	case "message-text":
		// Message of type text
	case "message-media":
		// Message with media attached
	default:
		return
	}
}
