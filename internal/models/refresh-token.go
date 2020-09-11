// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

// RefreshToken is used to validate a refresh token
type RefreshToken struct {
	Username string
	Token    string
}

func (r *RefreshToken) IsEmpty() bool {
	return r.Username != "" || r.Token != ""
}
