// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package store provides an interface for all database operations
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/constants"
	"chapper.dev/server/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Store wraps a GORM database connection
type Store struct {
	db       *gorm.DB
	settings *Settings
}

var (
	ErrInvalidDatabaseType = errors.New("Invalid database type")
)

type Settings struct {
	ID               uint `gorm:"primaryKey"`
	IsInstalled      bool
	SuperadminExists bool
}

// DefaultSettings provide the default values for settings
var DefaultSettings = &Settings{
	IsInstalled:      false,
	SuperadminExists: false,
}

// New returns a new MySQL DB instance
func New(t string, s config.StoreOptions) (*Store, error) {
	switch strings.ToLower(t) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(MySQLDSN(s)), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		db = db.Set("gorm:table_options", constants.StoreTableOptions)

		return &Store{
			db: db,
		}, nil
	default:
		return nil, ErrInvalidDatabaseType
	}
}

// MySQLDSN returns a data source name for a MySQL database connection, refer
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
func MySQLDSN(s config.StoreOptions) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		s.User,
		s.Password,
		s.Host,
		s.Port,
		s.Database,
		constants.StoreParams,
	)
}

// Migrate migrates the neccesary tables
func (s *Store) Migrate() error {
	return s.db.AutoMigrate(
		&models.Invite{},
		&models.Server{},
		&models.Room{},
		&models.User{},
		&models.Role{},
		&models.Privileges{},
		&Settings{},
	)
}

// GetSettings returns the settings or creates a new default entry in the database
func (s *Store) GetSettings() (*Settings, error) {
	settings := new(Settings)
	err := s.Ctx().Last(settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err := s.db.WithContext(context.Background()).Create(DefaultSettings).Error
			if err != nil {
				return nil, err
			}
			return DefaultSettings, nil
		}
		return nil, err
	}
	return settings, nil
}

// SetSettings sets settings and saves them
func (s *Store) SetSettings(settings *Settings) error {
	s.settings = settings
	return s.Ctx().Model(settings).Updates(settings).Error
}

// Ctx returns the db instance with a new context to create a new statement
func (s *Store) Ctx() *gorm.DB {
	// NOTE(Techassi): Apparently this is part of 'optimizations' in the new GORM version
	return s.db.WithContext(context.Background())
}
