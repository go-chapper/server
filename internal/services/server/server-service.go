// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package server provides utilities for creating, getting, updating and deleting virtual
// servers
package server

import (
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/store"
)

// Service wraps dependencies
type Service struct {
	store *store.Store
}

// NewService returns a new server service
func NewService(store *store.Store) Service {
	return Service{
		store: store,
	}
}

// CreateServer creates a new virtual server
func (s Service) CreateServer(server *models.Server) error {
	server.Hash = hash.Adler32(server.Name)
	return s.store.CreateServer(server)
}

// GetServer returns one virtual server identified by 'hash'
func (s Service) GetServer(hash string) (models.Server, error) {
	return s.store.GetServer(hash)
}

// GetServers returns all virtual servers
func (s Service) GetServers() ([]models.Server, error) {
	return s.store.GetServers()
}

// UpdateServer updates one virtual server identified by 'hash'
func (s Service) UpdateServer(hash string) error {
	return nil
}

// DeleteServer deletes one virtual server dentified by 'hash'
func (s Service) DeleteServer(hash string) error {
	return s.store.DeleteServer(hash)
}
