// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package user provides a user service to handle user related actions
package user

import (
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/store"
)

// Service is the top-level service struct
type Service struct {
	store *store.Store
}

// NewService returns a new user service
func NewService(store *store.Store) Service {
	return Service{
		store: store,
	}
}

// GetUser returns a user identified by the provided 'username' or an error if the user
// does not exist
func (s Service) GetUser(username string) (*models.User, error) {
	return s.store.GetUser(username)
}

// CreateUser creates a new 'user' or returns an error if the new user could not be
// created
func (s Service) CreateUser(user *models.User) error {
	// TODO <2020/10/09>: Can we optimize/improve this?
	settings, err := s.store.GetSettings()
	if err != nil {
		return err
	}

	if !settings.SuperadminExists {
		user.Role = append(user.Role, models.Superadmin())
		settings.SuperadminExists = true

		err := s.store.SetSettings(settings)
		if err != nil {
			return err
		}
	}

	return s.store.CreateUser(user)
}

func (s Service) SaveTempToken(username, token string) error {
	user, err := s.store.GetUser(username)
	if err != nil {
		return err
	}

	user.TwoFATempToken = token
	return s.store.SaveUser(user)
}
