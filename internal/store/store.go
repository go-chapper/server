// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"errors"
	"fmt"
	"strings"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/constants"
	"chapper.dev/server/internal/store/schemas"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

// Store wraps a sqlx database connection
type Store struct {
	conn     *sqlx.DB
	settings *Settings
}

var (
	// ErrInvalidDatabaseType indicates the provided database type is not supported
	ErrInvalidDatabaseType = errors.New("Invalid database type")
)

// Settings holds settings data
type Settings struct {
	ID               uint
	IsInstalled      bool
	SuperadminExists bool
}

// DefaultSettings provide the default values for settings
var DefaultSettings = &Settings{
	IsInstalled:      false,
	SuperadminExists: false,
}

// New returns a new store instance
func New(t string, options config.StoreOptions) (*Store, error) {
	switch strings.ToLower(t) {
	case "mysql":
		conn, err := sqlx.Open("mysql", DSN(options))
		if err != nil {
			return nil, err
		}

		return &Store{
			conn: conn,
		}, nil
	default:
		return nil, ErrInvalidDatabaseType
	}
}

// DSN returns a data source name for a database connection, refer
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
func DSN(options config.StoreOptions) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
		constants.StoreParams,
	)
}

// Migrate migrates the neccesary database tables
func (s *Store) Migrate() error {
	for _, scheme := range schemas.All() {
		_, err := s.conn.Exec(scheme)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSettings returns the settings or creates a new default entry in the database
func (s *Store) GetSettings() (*Settings, error) {
	if s.settings != nil {
		return s.settings, nil
	}
	// TODO <2020/10/12>: Re-implement this
	return nil, nil
}

// SetSettings sets settings and saves them
// func (s *Store) SetSettings(settings *Settings) error {
// 	s.settings = settings
// 	return s.Ctx().Model(settings).Updates(settings).Error
// }
