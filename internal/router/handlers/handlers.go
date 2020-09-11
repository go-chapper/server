// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/modules/jwt"
	"chapper.dev/server/internal/services/auth"
	"chapper.dev/server/internal/services/invite"
	"chapper.dev/server/internal/services/room"
	"chapper.dev/server/internal/services/server"
	"chapper.dev/server/internal/services/user"
	"chapper.dev/server/internal/store"

	j "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ErrorCode represents different error codes returned by the API
type ErrorCode string

const (
	ErrBind            ErrorCode = "bind-error"
	ErrEmptyData       ErrorCode = "empty-data"
	ErrInvalidData     ErrorCode = "invalid-data"
	ErrHashPassword    ErrorCode = "hash-error"
	ErrCreateUser      ErrorCode = "create-user-error"
	ErrUsernameTaken   ErrorCode = "username-taken"
	ErrServernameTaken ErrorCode = "servername-taken"
	ErrRoomnameTaken   ErrorCode = "roomname-taken"
	ErrGetUser         ErrorCode = "get-user-error"
	ErrUnauthorized    ErrorCode = "unauthorized"
	ErrTwoFA           ErrorCode = "2fa-error"
	ErrJWT             ErrorCode = "jwt-error"
	ErrCreateInvite    ErrorCode = "create-invite-error"
	ErrCreateServer    ErrorCode = "create-server-error"
	ErrCreateRoom      ErrorCode = "create-room-error"
)

type StatusCode string

const (
	StatusRegistered    StatusCode = "registered"
	StatusServerCreated StatusCode = "server-created"
	StatusRoomCreated   StatusCode = "room-created"
)

// NOTE(Techassi): Maybe split handlers for specific handler groups so we dont need to
// inject all services into one handler

// Handler provides an interface to handle different HTTP request
type Handler struct {
	config        *config.Config
	userService   user.Service
	authService   auth.Service
	inviteService invite.Service
	serverService server.Service
	roomService   room.Service
}

// Map is a wrapper for an map[string]interface{}, which gets used in JSON responses
type Map map[string]interface{}

// New returns a new handler with all required services injected
func New(store *store.Store, config *config.Config) *Handler {
	// Create services
	is := invite.NewService(store, *config)
	ss := server.NewService(store)
	rs := room.NewService(store)
	us := user.NewService(store)
	as := auth.NewService()

	return &Handler{
		config:        config,
		userService:   us,
		authService:   as,
		inviteService: is,
		serverService: ss,
		roomService:   rs,
	}
}

func getClaimes(c echo.Context) *jwt.Claims {
	user := c.Get("user").(*j.Token)
	return user.Claims.(*jwt.Claims)
}
