// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"time"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/services/errors"
	"chapper.dev/server/internal/store"
	"chapper.dev/server/internal/utils"

	"github.com/labstack/echo/v4"
)

var (
	// DefaultExpireTimespan describes the default expire timespan of an invite
	DefaultExpireTimespan = time.Hour * 24 * 7

	// inviteCtx describes the invite log context
	inviteCtx = log.NewContext("invite-srv")
)

// InviteService provides a service to create, get and delete invites
type InviteService struct {
	store  *store.Store
	config *config.Config
	logger *log.Logger
}

// NewInviteService returns a new invite service
func NewInviteService(store *store.Store, config *config.Config, logger *log.Logger) InviteService {
	return InviteService{
		store:  store,
		config: config,
		logger: logger,
	}
}

// CreateInvite creates a new invite link
func (s InviteService) CreateInvite(username string, c echo.Context) (*models.Invite, error) {
	var invite *models.Invite

	// Bind to invite model
	err := c.Bind(invite)
	if err != nil {
		s.logger.Errorc(inviteCtx, err)
		return nil, errors.ErrBindInvite
	}

	// Check if some data is missing
	if invite.IsEmpty() {
		s.logger.Infoc(inviteCtx, "some data to create an invite is missing")
		return nil, errors.ErrMissingInviteData
	}

	// Get expire time and set invite values
	expireTime := time.Now().Add(DefaultExpireTimespan)

	invite.CreatedBy = username
	invite.ExpiresAt = utils.ToNullTime(expireTime)
	invite.Hash = hash.Adler32(invite.Server + expireTime.String())

	// Finally store the invite to the database
	err = s.store.CreateInvite(invite)
	if err != nil {
		s.logger.Errorc(inviteCtx, err)
		return nil, errors.ErrCreateInvite
	}

	return invite, nil
}
