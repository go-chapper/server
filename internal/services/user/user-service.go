// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package user provides a user service to handle user related actions
package user

import (
	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/models/joins"
	"chapper.dev/server/internal/modules/avatar"
	"chapper.dev/server/internal/store"
)

// Service is the top-level service struct
type Service struct {
	store  *store.Store
	config config.Config
}

// NewService returns a new user service
func NewService(s *store.Store, c config.Config) Service {
	return Service{
		store:  s,
		config: c,
	}
}

// GetUser returns a user identified by the provided 'username' or an error if the user
// does not exist
func (s Service) GetUser(username string) (*models.User, error) {
	return s.store.GetUser(username)
}

func (s Service) GetUserPublicKey(username string) (string, error) {
	return s.store.GetUserPublicKey(username)
}

// CreateUser creates a new 'user' or returns an error if the new user could not be
// created
func (s Service) CreateUser(newUser *models.SignupUser) error {
	// TODO <2020/10/09>: Can we optimize/improve this?
	settings, err := s.store.GetSettings()
	if err != nil {
		return err
	}

	user := &models.User{
		Username:  newUser.Username,
		Password:  newUser.Password,
		Email:     newUser.Email,
		PublicKey: newUser.PublicKey,
	}

	if !settings.SuperadminExists {
		user.Role = append(user.Role, models.Superadmin())
		settings.SuperadminExists = true

		err := s.store.SetSettings(settings)
		if err != nil {
			return err
		}
	} else {
		user.Role = append(user.Role, models.Basic())
	}

	err = s.store.CreateUser(user)
	if err != nil {
		return err
	}

	a := avatar.New(240, user.Username)
	return a.Generate(s.config.Router.AvatarPath)
}

func (s Service) SaveTempToken(username, token string) error {
	user, err := s.store.GetUser(username)
	if err != nil {
		return err
	}

	user.TwoFATempToken = token
	return s.store.SaveUser(user)
}

func (s Service) GetUserServers(username string) ([]joins.UserServers, error) {
	return s.store.GetUserServers(username)
}

func (s Service) PutUserServer(username string) ([]models.Server, error) {
	return nil, nil
}
