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
	"chapper.dev/server/internal/transport/signaling"

	j "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ErrorCode represents different error codes returned by the API
type ErrorCode string

const (
	ErrInternal        ErrorCode = "internal-error"
	ErrBind            ErrorCode = "bind-error"
	ErrEmptyData       ErrorCode = "empty-data"
	ErrInvalidData     ErrorCode = "invalid-data"
	ErrGetUser         ErrorCode = "get-user-error"
	ErrHashPassword    ErrorCode = "hash-error"
	ErrCreateUser      ErrorCode = "create-user-error"
	ErrUsernameTaken   ErrorCode = "username-taken"
	ErrUnauthorized    ErrorCode = "unauthorized"
	ErrTwoFA           ErrorCode = "2fa-error"
	ErrJWT             ErrorCode = "jwt-error"
	ErrServernameTaken ErrorCode = "servername-taken"
	ErrCreateServer    ErrorCode = "create-server-error"
	ErrServerNotFound  ErrorCode = "server-not-found"
	ErrRoomnameTaken   ErrorCode = "roomname-taken"
	ErrCreateRoom      ErrorCode = "create-room-error"
	ErrRoomNotFound    ErrorCode = "room-not-found"
	ErrCreateInvite    ErrorCode = "create-invite-error"
	ErrCreateWebsocket ErrorCode = "create-websocket-error"
	ErrCreateDirect    ErrorCode = "create-direct-error"
)

type StatusCode string

const (
	StatusRegistered    StatusCode = "registered"
	StatusServerCreated StatusCode = "server-created"
	StatusRoomCreated   StatusCode = "room-created"
	StatusServerDeleted StatusCode = "server-deleted"
	StatusRoomDeleted   StatusCode = "room-deleted"
)

// NOTE(Techassi): Maybe split handlers for specific handler groups so we dont need to
// inject all services into one handler

// Handler provides an interface to handle different HTTP request
type Handler struct {
	config        *config.Config
	hub           *signaling.Hub
	userService   user.Service
	authService   auth.Service
	inviteService invite.Service
	serverService server.Service
	roomService   room.Service
}

// Map is a wrapper for an map[string]interface{}, which gets used in JSON responses
type Map map[string]interface{}

// New returns a new handler with all required services injected
func New(store *store.Store, config *config.Config, hub *signaling.Hub) *Handler {
	// Create services
	is := invite.NewService(store, *config)
	us := user.NewService(store, *config)
	ss := server.NewService(store)
	rs := room.NewService(store)
	as := auth.NewService()

	return &Handler{
		config:        config,
		hub:           hub,
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
