// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

import "encoding/json"

// Typed is a generic typed message
type Typed struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Message defines an interface for a more specific message
type Message interface {
	// Handle handles the execution of message related tasks. If the tasks fail, an error
	// is returned
	Handle(*Hub) error

	// Type returns the type of the message as a string
	Type() string

	// ToTyped returns the message as a generic typed message
	ToTyped() (Typed, error)

	// New returns a constructor to initialize a new message
	New() func() Message
}

func AllMessages() []Message {
	return []Message{&AvailabilityChange{}}
}

// ToMessage returns the generic typed message as a specific message
func (t Typed) ToMessage(h *Hub) (Message, error) {
	init, exists := h.messages[t.Type]
	if !exists {
		return nil, ErrInvalidMessageType
	}

	message := init()
	err := json.Unmarshal(t.Data, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
