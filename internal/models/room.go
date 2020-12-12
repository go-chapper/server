// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import "gopkg.in/guregu/null.v4"

type Room struct {
	Hash        string      `json:"hash" db:"hash"`
	Name        string      `json:"name" db:"name"`
	Type        null.String `json:"type" db:"type"`
	Description null.String `json:"description" db:"description"`
}

var allowedTypes = []string{"text", "voice"}

// IsEmpty returns if all required data is present
func (r *Room) IsEmpty() bool {
	return r.Name == "" || r.Type.String == ""
}

// Invalid returns if the data is invalid
func (r *Room) Invalid() bool {
	if !isIn(r.Type.String, allowedTypes) {
		return true
	}

	return false
}

func isIn(t string, types []string) bool {
	for i := 0; i < len(types); i++ {
		if t == types[i] {
			return true
		}
	}
	return false
}
