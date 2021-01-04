// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

// func (h *Hub) handleRegister(c *Connection) {
// 	h.Lock()
// 	defer h.Unlock()

// 	h.conns[c] = ""
// }

// func (h *Hub) handleUnregister(c *Connection) {
// 	h.Lock()
// 	defer h.Unlock()

// 	username, ok := h.conns[c]
// 	if !ok {
// 		return
// 	}

// 	c.Close()
// 	delete(h.conns, c)
// 	delete(h.peers, username)
// }

func (h *Hub) handleMessage(m *Message) {
	// switch m.Type {
	// case "subscribe":
	// 	if m.Token == "" || m.Username == "" {
	// 		return
	// 	}

	// 	t, ok := h.tokens[m.Username]
	// 	if !ok || t != m.Token {
	// 		go h.Unregister(m.connection)
	// 		return
	// 	}

	// 	h.peers[m.Username] = &Peer{
	// 		Username:   m.Username,
	// 		Token:      m.Token,
	// 		connection: m.connection,
	// 	}
	// 	h.conns[m.connection] = m.Username
	// case "status-writing":
	// 	// Indicate a chat member is writing
	// case "message-text":
	// 	// Message of type text
	// 	if len(m.To) == 0 || m.From == "" {
	// 		return
	// 	}

	// 	h.Send(m)
	// case "message-media":
	// 	// Message with media attached
	// default:
	// 	return
	// }
}
