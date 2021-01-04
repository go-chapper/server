// Copyright (c) 2021-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

type AuthenticationMessage struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// Handle handles the authentication of a user
func (a *AuthenticationMessage) Handle(h *Hub) error {
	return nil
}

// Type returns the type of this message as a string
func (a *AuthenticationMessage) Type() string {
	return "authentication"
}

// New returns a function to create a new AuthenticationMessage
func (a *AuthenticationMessage) New() func() Message {
	return func() Message {
		return a
	}
}
