// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"chapper.dev/server/internal/models"
)

// CreateRoom inserts a new room entry into the database
func (s *Store) CreateRoom(room *models.Room) error {
	_, err := s.conn.Exec(`
		INSERT INTO rooms
		(hash, name, type, description)
		VALUES (?, ?, ?, ?)`,
		room.Hash,
		room.Name,
		room.Type,
		room.Description,
	)
	return err
}

// GetRoom selects ONE room entry with provided 'roomHash' from the database
func (s *Store) GetRoom(roomHash string) (models.Room, error) {
	var room models.Room
	err := s.conn.Get(&room,
		`SELECT hash, name, type, description 
		FROM rooms 
		WHERE hash = ?`,
		roomHash,
	)
	return room, err
}

// GetRooms selects multiple room entries from the database
func (s *Store) GetRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := s.conn.Select(&rooms, `SELECT hash, name, type, description FROM rooms`)
	return rooms, err
}

// UpdateRoom updates ONE room entry with provided 'roomHash' in the database
func (s *Store) UpdateRoom(roomHash string, new *models.Room) error {
	_, err := s.conn.Exec(`
		UPDATE rooms
		SET name = ?, type = ?, description = ?
		WHERE hash = ?`,
		new.Name,
		new.Type,
		new.Description,
		roomHash,
	)
	return err
}

// DeleteRoom deletes ONE room entry with provided 'roomHash' from the database
func (s *Store) DeleteRoom(roomHash string) error {
	_, err := s.conn.Exec(`
		DELETE FROM rooms
		WHERE hash = ?`,
		roomHash,
	)
	return err
}
