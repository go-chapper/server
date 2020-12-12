// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package app provides the central entry point to read in the config, setup the server
// and start services
package app

import (
	"fmt"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/router"
	"chapper.dev/server/internal/router/handlers"
	"chapper.dev/server/internal/store"
	"chapper.dev/server/internal/transport/turn"
)

var appCtx = log.NewContext("app")

type App struct {
	config *config.Config
	logger *log.Logger
	store  *store.Store
	router *router.Router
	turn   *turn.TURN
}

func New(configFilePath string) (*App, error) {
	cfg := config.New()
	err := cfg.Read(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("[E] [app] failed to read config file: %v", err)
	}

	logger, err := log.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("[E] [app] failed to setup logger: %v", err)
	}

	db, err := store.New("mysql", cfg.Store)
	if err != nil {
		logger.Errorc(appCtx, err)
		return nil, err
	}

	err = db.Migrate()
	if err != nil {
		logger.Errorc(appCtx, err)
		return nil, err
	}

	rauter := router.New(cfg)
	handle := handlers.New(db, cfg, logger)
	rauter.AddRoutes(handle)

	turnServer, err := turn.New(cfg.Turn.PublicIP, cfg.Router.Domain, "udp4", cfg.Turn.Port)
	if err != nil {
		logger.Errorc(appCtx, err)
		return nil, err
	}

	return &App{
		config: cfg,
		logger: logger,
		store:  db,
		router: rauter,
		turn:   turnServer,
	}, nil
}

func (a *App) Run() error {
	err := a.turn.Run()
	if err != nil {
		a.logger.Errorc(appCtx, err)
		return err
	}

	err = a.router.Run()
	if err != nil {
		a.logger.Errorc(appCtx, err)
		return err
	}

	return nil
}
