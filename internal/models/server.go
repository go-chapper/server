// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

type Server struct {
	Hash        string `json:"hash" gorm:"primaryKey"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	// Rooms       []Room `json:"rooms,omitepmty" gorm:"foreignKey:BelongsTo"`
}

func (s *Server) IsEmpty() bool {
	return s.Name == ""
}
