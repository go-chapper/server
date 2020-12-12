// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
)

type Invite struct {
	Hash       string    `json:"hash" db:"hash"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	Server     string    `json:"server" db:"server"`
	OneTimeUse bool      `json:"one_time_use" db:"one_time_use"`
	ExpiresAt  null.Time `json:"expires_at" db:"expires_at"`
}

// ToURL returns the URL representation of the invite
func (i *Invite) ToURL(domain string) string {
	return fmt.Sprintf("https://%s/i/%s", domain, i.Hash)
}

// IsEmpty returns if all data is present
func (c *Invite) IsEmpty() bool {
	return c.Server == "" || (!c.OneTimeUse && c.ExpiresAt.IsZero())
}
