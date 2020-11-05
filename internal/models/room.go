// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

type Room struct {
	Hash        string `json:"hash" gorm:"primaryKey"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	BelongsTo   string `json:"-" gorm:"index"` // Reference
}

var allowedTypes = []string{"text", "voice"}

// IsEmpty returns if all required data is present
func (r *Room) IsEmpty() bool {
	return r.Name == "" || r.Type == ""
}

// Invalid returns if the data is invalid
func (r *Room) Invalid() bool {
	if !isIn(r.Type, allowedTypes) {
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
