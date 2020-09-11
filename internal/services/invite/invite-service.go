// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package invite provides utilities to create and validate invite links
package invite

import (
	"fmt"
	"net/url"
	"time"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/store"
)

type Service struct {
	store  *store.Store
	config config.Config
}

func NewService(store *store.Store, config config.Config) Service {
	return Service{
		store:  store,
		config: config,
	}
}

// CreateInvite creates a new invite link
func (s Service) CreateInvite(createdBy, server string, oneTimeUse bool) (*models.Invite, error) {
	var (
		currentTime = time.Now()
		hash        = hash.Adler32(server + currentTime.String())
	)

	// TODO <2020/07/09>: Dont hardcode http
	url, err := url.Parse(fmt.Sprintf("http://%s/i/%s", s.config.Router.Domain, hash))
	if err != nil {
		return nil, err
	}

	invite := &models.Invite{
		CreatedBy:  createdBy,
		Hash:       hash,
		Server:     server,
		URL:        url,
		URLString:  url.String(),
		OneTimeUse: oneTimeUse,
		ExpiresAt:  currentTime.Add(time.Hour),
	}

	err = s.store.CreateInvite(invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}
