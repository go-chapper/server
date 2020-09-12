// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

type Direct struct {
	Hash string `json:"hash" gorm:"primaryKey"`
	Name string `json:"name"`
}
