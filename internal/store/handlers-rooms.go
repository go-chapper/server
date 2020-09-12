// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"chapper.dev/server/internal/models"
)

// CreateRoom creates a new 'room' on server identified by 'serverHash'
func (s *Store) CreateRoom(serverHash string, room *models.Room) error {
	return s.Ctx().Model(&models.Server{Hash: serverHash}).Association("Rooms").Append(room)
}

// GetRoom returns one room of a server identified by 'serverHash'
func (s *Store) GetRoom(serverHash, roomHash string) (models.Room, error) {
	var room models.Room
	return room, s.Ctx().Model(&models.Server{Hash: serverHash}).Association("Rooms").Find(&room, "hash = ?", roomHash)
}

// GetRooms returns all rooms of a server identified by 'serverHash'
func (s *Store) GetRooms(serverHash string) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, s.Ctx().Model(&models.Server{Hash: serverHash}).Association("Rooms").Find(&rooms)
}

// UpdateRoom updates a room of a server identified by 'serverHash'
func (s *Store) UpdateRoom(serverHash, roomHash string) error {
	return nil
}

// DeleteRoom deletes a room of a server identified by 'serverHash'
func (s *Store) DeleteRoom(serverHash, roomHash string) error {
	return s.Ctx().Model(&models.Server{Hash: serverHash}).Association("Rooms").Delete(&models.Room{Hash: roomHash})
}
