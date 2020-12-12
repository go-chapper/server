// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import "gopkg.in/guregu/null.v4"

type Server struct {
	Hash        string      `json:"hash" db:"hash"`
	Name        string      `json:"name" db:"name"`
	Description null.String `json:"description" db:"description"`
	Image       null.String `json:"image" db:"image"`
}

func (s *Server) IsEmpty() bool {
	return s.Name == ""
}
