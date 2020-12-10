// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

type Server struct {
	Hash        string `json:"hash" db:"hash"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Image       string `json:"image" db:"image"`
}

func (s *Server) IsEmpty() bool {
	return s.Name == ""
}
