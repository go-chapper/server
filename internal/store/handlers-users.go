// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"fmt"

	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/models/joins"
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

func (s *Store) GetUserServers(username string) ([]joins.UserServers, error) {
	servers := []joins.UserServers{}
	err := s.db.Raw(`
		SELECT 
			servers.* 
		FROM 
			servers 
		LEFT JOIN (users, user_servers) 
			ON (
				user_servers.user_username = users.username AND 
				user_servers.server_hash = servers.hash
			) 
		WHERE users.username = ?`, username).
		Scan(&servers).Error
	fmt.Println(servers, username)
	return servers, err
}

func (s *Store) CreateUser(user *models.User) error {
	return s.Ctx().Create(user).Error
}

func (s *Store) SaveUser(user *models.User) error {
	return s.Ctx().Save(user).Error
}
