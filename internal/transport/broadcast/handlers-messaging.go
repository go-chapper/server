// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package broadcast provides utilities to broadcast messages
package broadcast

import "fmt"

func (h *MessagingHub) handleRegister(c *Connection) {
	h.Lock()
	defer h.Unlock()

	h.conns[c] = true
}

func (h *MessagingHub) handleUnregister(c *Connection) {
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

func (h *MessagingHub) handleMessage(m *Message) {

}
