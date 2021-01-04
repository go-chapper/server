// Copyright (c) 2021-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

type TextMessage struct {
}

// Handle handles the forwarding of a text message
func (t *TextMessage) Handle(h *Hub) error {
	return nil
}

// Type returns the type of this message as a string
func (t *TextMessage) Type() string {
	return "authentication"
}

// New returns a function to create a new TextMessage
func (t *TextMessage) New() func() Message {
	return func() Message {
		return t
	}
}
