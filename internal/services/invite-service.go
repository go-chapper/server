// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"errors"
	"time"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/store"
)

var (
	ErrMissingData = errors.New("invite-missing-data")
)

// InviteService wrapps all dependencies of the invite service
type InviteService struct {
	store  *store.Store
	config *config.Config
}

// NewService returns a new invite service
func NewInviteService(store *store.Store, config *config.Config) InviteService {
	return InviteService{
		store:  store,
		config: config,
	}
}

// CreateInvite creates a new invite link
func (s InviteService) CreateInvite(createdBy string, newInvite *models.CreateInvite) (*models.Invite, error) {
	if newInvite.IsEmpty() {
		return nil, ErrMissingData
	}

	var (
		currentTime = time.Now()
		hash        = hash.Adler32(newInvite.Server + currentTime.String())
	)

	invite := &models.Invite{
		CreatedBy:  createdBy,
		Hash:       hash,
		Server:     newInvite.Server,
		OneTimeUse: newInvite.OneTimeUse,
		ExpiresAt:  newInvite.ExpiresAt,
	}

	err := s.store.CreateInvite(invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}
