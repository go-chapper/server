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
	Host       string    `json:"host"`
	Server     string    `json:"server"`
	URL        *url.URL  `gorm:"-"`
	URLString  string    `json:"-"`
	OneTimeUse bool      `json:"oneTimeUse"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// CreateInvite binds to the request to create a new invite
type CreateInvite struct {
	Host       string    `json:"host"`
	Server     string    `json:"server"`
	OneTimeUse bool      `json:"oneTimeUse"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// IsEmpty returns if all data is present
func (c *CreateInvite) IsEmpty() bool {
	return c.Host == "" || c.Server == "" || (!c.OneTimeUse && c.ExpiresAt.IsZero())
}
