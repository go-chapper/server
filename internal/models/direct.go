// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

type Direct struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
	// SessionDescription string `json:"session_description"`
}

func (d *Direct) IsEmpty() bool {
	return d.Caller == "" || d.Callee == ""
}
