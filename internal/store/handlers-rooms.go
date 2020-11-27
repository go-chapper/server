// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"chapper.dev/server/internal/models"
)

// CreateRoom creates a new 'room'
func (s *Store) CreateRoom(room *models.Room) error {
	return s.Ctx().Create(room).Error
}

// GetRoom returns one room indentified by 'roomHash'
func (s *Store) GetRoom(roomHash string) (models.Room, error) {
	var room models.Room
	return room, s.Ctx().Find(&room, "hash = ?", roomHash).Error
}

// GetRooms returns all rooms
func (s *Store) GetRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, s.Ctx().Find(&rooms).Error
}

// UpdateRoom updates a room identified by 'roomHash'
func (s *Store) UpdateRoom(roomHash string) error {
	return nil
}

// DeleteRoom deletes a room identified by 'roomHash'
func (s *Store) DeleteRoom(roomHash string) error {
	return s.Ctx().Where("hash = ?", roomHash).Delete(&models.Room{}).Error
}
