// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"chapper.dev/server/internal/models"
)

func (s *Store) GetUser(username string) (*models.User, error) {
	user := new(models.User)
	return user, s.Ctx().
		Preload("Role.Privileges").
		Where("username = ?", username).
		First(user).Error
}

func (s *Store) GetUserPublicKey(username string) (string, error) {
	user := new(models.User)
	return user.PublicKey,
		s.Ctx().
			Where("username = ?", username).
			Select("public_key").
			First(user).Error
}

func (s *Store) GetUserServers(username string) ([]models.Server, error) {
	servers := []models.Server{}
	err := s.Ctx().
		Model(models.User{Username: username}).
		Association("Servers").
		Find(&servers)
	return servers, err
}

func (s *Store) CreateUser(user *models.User) error {
	return s.Ctx().Create(user).Error
}

func (s *Store) SaveUser(user *models.User) error {
	return s.Ctx().Save(user).Error
}
