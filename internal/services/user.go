// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/avatar"
	"chapper.dev/server/internal/store"
)

// UserService is the top-level service struct
type UserService struct {
	store  *store.Store
	config *config.Config
}

// NewUserService returns a new user service
func NewUserService(s *store.Store, c *config.Config) UserService {
	return UserService{
		store:  s,
		config: c,
	}
}

// GetUser returns a user identified by the provided 'username' or an error if the user
// does not exist
func (s UserService) GetUser(username string) (models.User, error) {
	return s.store.GetUser(username)
}

func (s UserService) GetUserPublicKey(username string) (string, error) {
	return s.store.GetUserPublicKey(username)
}

// CreateUser creates a new 'user' or returns an error if the new user could not be
// created
func (s UserService) CreateUser(user models.PublicUser) error {
	// TODO <2020/10/09>: Can we optimize/improve this?
	// settings, err := s.store.GetSettings()
	// if err != nil {
	// 	return err
	// }

	// if !settings.SuperadminExists {
	// 	user.Role = append(user.Role, models.Superadmin())
	// 	settings.SuperadminExists = true

	// 	err := s.store.SetSettings(settings)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	user.Role = append(user.Role, models.Basic())
	// }

	err := s.store.CreateUser(user)
	if err != nil {
		return err
	}

	a := avatar.New(240, user.Username)
	return a.Generate(s.config.Router.AvatarPath)
}

func (s UserService) UpdateTwoFAVerify(username, verify string) error {
	return s.store.UpdateTwoFAVerify(username, verify)
}

func (s UserService) PutUserServer(username string) ([]models.Server, error) {
	return nil, nil
}
