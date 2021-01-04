// Copyright (c) 2021-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

import "chapper.dev/server/internal/constants"

type TypingChange struct {
	Scope    string                `json:"scope"`
	Username string                `json:"username"`
	State    constants.TypingState `json:"state"`
}

// Handle handles the change of the typing state of one user
func (t *TypingChange) Handle(h *Hub) error {
	return nil
}

// Type returns the type of this message as a string
func (t *TypingChange) Type() string {
	return "typing-change"
}

// New returns a function to create a new TypingChange message
func (t *TypingChange) New() func() Message {
	return func() Message {
		return t
	}
}
