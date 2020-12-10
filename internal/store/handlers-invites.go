// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"chapper.dev/server/internal/models"
)

// CreateInvite creates a new invite
func (s *Store) CreateInvite(invite *models.Invite) error {
	// TODO <2020/10/12>: set created_at & expires_at timestamp
	_, err := s.conn.Exec(`
		INSERT INTO invites
		(hash, created_by, server, one_time_use, expires_at)
		VALUES (?, ?, ?, ?, ?)`,
		invite.Hash,
		invite.CreatedBy,
		invite.Server,
		invite.OneTimeUse,
		invite.ExpiresAt,
	)
	return err
}
