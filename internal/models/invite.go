// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"time"
)

type Invite struct {
	Hash       string    `json:"hash" db:"hash"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	Server     string    `json:"server" db:"server"`
	OneTimeUse bool      `json:"one_time_use" db:"one_time_use"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
}

// CreateInvite binds to the request to create a new invite
type CreateInvite struct {
	Server     string    `json:"server"`
	OneTimeUse bool      `json:"oneTimeUse"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// ToURL returns the URL representation of the invite
func (i *Invite) ToURL(domain string) string {
	return fmt.Sprintf("https://%s/i/%s", domain, i.Hash)
}

// IsEmpty returns if all data is present
func (c *CreateInvite) IsEmpty() bool {
	return c.Server == "" || (!c.OneTimeUse && c.ExpiresAt.IsZero())
}
