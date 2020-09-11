// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

import (
	"net/url"
	"time"
)

type Invite struct {
	Hash       string    `json:"-" gorm:"primaryKey"`
	CreatedBy  string    `json:"-"`
	Server     string    `json:"server"`
	URL        *url.URL  `gorm:"-"`
	URLString  string    `json:"-"`
	OneTimeUse bool      `json:"oneTimeUse"`
	ExpiresAt  time.Time `josn:"expiresAt"`
}
