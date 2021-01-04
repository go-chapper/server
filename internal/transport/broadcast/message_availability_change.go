// Copyright (c) 2021-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package broadcast

import "chapper.dev/server/internal/constants"

// AvailabilityChange defines the event when the availability status of a user changes.
type AvailabilityChange struct {
	Username string                      `json:"username"`
	State    constants.AvailabilityState `json:"state"`
}

// Handle handles the change of the availability change of one user
func (a *AvailabilityChange) Handle(h *Hub) error {
	return nil
}

// Type returns the type of this message as a string
func (a *AvailabilityChange) Type() string {
	// NOTE(Techassi): Should we define this in the interface or should we just define this once on the hub?
	return "availability-change"
}

// ToTyped returns this message as a generic typed message
func (a *AvailabilityChange) ToTyped() (Typed, error) {
	return Typed{}, nil
}

// New returns a function to create a new AvailabilityChange message
func (a *AvailabilityChange) New() func() Message {
	return func() Message {
		return a
	}
}
