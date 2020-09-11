// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"git.web-warrior.de/go-chapper/server/internal/models"
)

// CreateInvite creates a new invite
func (s *Store) CreateInvite(invite *models.Invite) error {
	return s.Ctx().Create(invite).Error
}
