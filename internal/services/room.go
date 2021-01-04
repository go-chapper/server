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

var roomCtx = log.NewContext("room-srv")

// RoomService provides a service to create, get, update and delete rooms
type RoomService struct {
	store  *store.Store
	logger *log.Logger
}

// NewRoomService returns a new room service
func NewRoomService(store *store.Store, logger *log.Logger) RoomService {
	return RoomService{
		store:  store,
		logger: logger,
	}
}

// CreateRoom creates and inserts a new room into the database
func (s RoomService) CreateRoom(c echo.Context) error {
	var room = new(models.Room)

	// Bind to room model
	err := c.Bind(room)
	if err != nil {
		s.logger.Errorc(roomCtx, err)
		return errors.ErrBindRoom
	}

	// Check if some data is missing or invalid
	if room.IsEmpty() || room.Invalid() {
		s.logger.Infoc(roomCtx, "data missing to create room")
		return errors.ErrMissingRoomData
	}

	// Calculate room hash and insert into database
	room.Hash = hash.FNV64(room.Name)
	err = s.store.CreateRoom(room)
	if err != nil {
		s.logger.Errorc(roomCtx, err)
		return errors.ErrCreateRoom
	}

	return nil
}

// GetRoom returns ONE room from the database identified by the provided hash
func (s RoomService) GetRoom(c echo.Context) (*models.Room, error) {
	roomHash := c.Param("room-hash")

	if roomHash == "" {
		s.logger.Infoc(roomCtx, "invalid room hash")
		return nil, errors.ErrInvalidHash
	}

	return s.store.GetRoom(roomHash)
}

// GetRooms returns multiple rooms from the database
func (s RoomService) GetRooms() ([]models.Room, error) {
	return s.store.GetRooms()
}

// UpdateRoom updates ONE room in the database indentified by the provided hash
func (s RoomService) UpdateRoom(c echo.Context) error {
	var newRoom *models.Room
	roomHash := c.Param("room-hash")

	err := c.Bind(newRoom)
	if err != nil {
		s.logger.Errorc(roomCtx, err)
		return errors.ErrBindRoom
	}

	err = s.store.UpdateRoom(roomHash, newRoom)
	if err != nil {
		s.logger.Errorc(roomCtx, err)
		return errors.ErrUpdateRoom
	}

	return nil
}

// DeleteRoom deletes ONE room from the database identified by the provided room hash
func (s RoomService) DeleteRoom(c echo.Context) error {
	roomHash := c.Param("room-hash")
	return s.store.DeleteRoom(roomHash)
}
