// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package messages provides in- and outgoing WebRTC messages
package messages

import "chapper.dev/server/internal/transport/signaling"

// CallOffer is a call offer send from the sender to receiver with a session description
// to start a call
type CallOffer struct {
	Sender   string
	Receiver string
	SDP      string
}

func (c *CallOffer) Handle(h *signaling.Hub) error {
	_ = h.GetClient(c.Receiver)
	return nil
}
