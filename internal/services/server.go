// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/store"
)

// ServerService wraps dependencies
type ServerService struct {
	store *store.Store
}

// NewServerService returns a new server service
func NewServerService(store *store.Store) ServerService {
	return ServerService{
		store: store,
	}
}

// CreateServer creates a new virtual server
func (s ServerService) CreateServer(server *models.Server) error {
	server.Hash = hash.Adler32(server.Name)
	return s.store.CreateServer(server)
}

// GetServer returns one virtual server identified by 'hash'
func (s ServerService) GetServer(hash string) (models.Server, error) {
	return s.store.GetServer(hash)
}

// GetServers returns all virtual servers
func (s ServerService) GetServers() ([]models.Server, error) {
	return s.store.GetServers()
}

// UpdateServer updates one virtual server identified by 'hash'
func (s ServerService) UpdateServer(hash string) error {
	return nil
}

// DeleteServer deletes one virtual server dentified by 'hash'
func (s ServerService) DeleteServer(hash string) error {
	return s.store.DeleteServer(hash)
}
