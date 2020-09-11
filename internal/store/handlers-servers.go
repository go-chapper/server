// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"chapper.dev/server/internal/models"
)

func (s *Store) CreateServer(server *models.Server) error {
	return s.Ctx().Create(server).Error
}

func (s *Store) GetServer(hash string) (models.Server, error) {
	var server models.Server
	return server, s.Ctx().Preload("Rooms").Where("hash = ?", hash).First(&server).Error
}

func (s *Store) GetServers() ([]models.Server, error) {
	var servers []models.Server
	return servers, s.Ctx().Find(&servers).Error
}

func (s *Store) UpdateServer() error {
	return nil
}

func (s *Store) DeleteServer(hash string) error {
	return s.Ctx().Where("hash = ?", hash).Delete(&models.Server{}).Error
}
