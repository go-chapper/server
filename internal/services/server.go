// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package services

import (
	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/models"
	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/services/errors"
	"chapper.dev/server/internal/store"

	"github.com/labstack/echo/v4"
)

var serverCtx = log.NewContext("server-srv")

// ServerService wraps dependencies
type ServerService struct {
	store  *store.Store
	logger *log.Logger
}

// NewServerService returns a new server service
func NewServerService(store *store.Store, logger *log.Logger) ServerService {
	return ServerService{
		store:  store,
		logger: logger,
	}
}

// CreateServer creates a new virtual server
func (s ServerService) CreateServer(c echo.Context) error {
	var server *models.Server

	err := c.Bind(server)
	if err != nil {
		s.logger.Errorc(serverCtx, err)
		return errors.ErrBindServer
	}

	if server.IsEmpty() {
		s.logger.Infoc(serverCtx, "data missing to create server")
		return errors.ErrMissingServerData
	}

	serverHash := hash.FNV64(server.Name)
	server.Hash = serverHash

	err = s.store.CreateServer(server)
	if err != nil {
		s.logger.Errorc(serverCtx, err)
		return errors.ErrCreateServer
	}

	return nil
}

// GetServer returns one virtual server identified by 'hash'
func (s ServerService) GetServer(c echo.Context) (*models.Server, error) {
	serverHash := c.Param("server-hash")

	if serverHash == "" {
		s.logger.Infoc(serverCtx, "invalid server hash")
		return nil, errors.ErrInvalidHash
	}

	return s.store.GetServer(serverHash)
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
