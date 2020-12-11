// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"chapper.dev/server/internal/models"
)

// CreateInvite creates a new invite
func (s *Store) CreateInvite(invite *models.Invite) error {
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
