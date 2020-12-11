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
	var room *models.Room

	err := c.Bind(room)
	if err != nil {
		s.logger.Errorc(roomCtx, err)
		return errors.ErrBindRoom
	}

	if room.IsEmpty() || room.Invalid() {
		s.logger.Infoc(roomCtx, "data missing to create room")
		return errors.ErrMissingRoomData
	}

	room.Hash = hash.FNV64(room.Name)
	err = s.store.CreateRoom(room)
	if err != nil {
		s.logger.Errorc(roomCtx, err)
		return errors.ErrCreateRoom
	}

	return nil
}

func (s RoomService) GetRoom(serverHash string) (models.Room, error) {
	return s.store.GetRoom(serverHash)
}

func (s RoomService) GetRooms() ([]models.Room, error) {
	return s.store.GetRooms()
}

func (s RoomService) UpdateRoom(roomHash string) error {
	return nil
}

func (s RoomService) DeleteRoom(roomHash string) error {
	return s.store.DeleteRoom(roomHash)
}
