// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"net/http"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/log"
	"chapper.dev/server/internal/modules/jwt"
	"chapper.dev/server/internal/services"
	"chapper.dev/server/internal/services/errors"
	"chapper.dev/server/internal/store"

	j "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ErrorCode represents different error codes returned by the API
type ErrorCode string

const (
	ErrInternal     ErrorCode = "internal-error"
	ErrUnauthorized ErrorCode = "unauthorized"
)

var handlerCtx = log.NewContext("handler")

// Handler provides an interface to handle different HTTP request
type Handler struct {
	config *config.Config
	logger *log.Logger
	// signalingHub  broadcast.Hub
	// messagingHub  broadcast.Hub
	inviteService services.InviteService
	serverService services.ServerService
	userService   services.UserService
	authService   services.AuthService
	roomService   services.RoomService
	callService   services.CallService
}

// Map is a wrapper for an map[string]interface{}, which gets used in JSON responses
type Map map[string]interface{}

// New returns a new handler with all required services injected
func New(store *store.Store, config *config.Config, logger *log.Logger) *Handler {
	// Create services
	is := services.NewInviteService(store, config, logger)
	as := services.NewAuthService(store, config, logger)
	ss := services.NewServerService(store, logger)
	us := services.NewUserService(store, config)
	rs := services.NewRoomService(store, logger)
	cs := services.NewCallService()

	// signalingHub := broadcast.NewSignalingHub()
	// messagingHub := broadcast.NewMessagingHub()

	return &Handler{
		config: config,
		logger: logger,
		// signalingHub:  signalingHub,
		// messagingHub:  messagingHub,
		userService:   us,
		authService:   as,
		inviteService: is,
		serverService: ss,
		roomService:   rs,
		callService:   cs,
	}
}

func (h *Handler) handleError(err error, c echo.Context) error {
	if se, ok := err.(*errors.ServiceError); ok {
		h.logger.Errorc(handlerCtx, se)
		return c.JSON(se.Code(), Map{
			"error": se.Err(),
		})
	}

	h.logger.Errorc(handlerCtx, err)
	return c.JSON(http.StatusInternalServerError, Map{
		"error": ErrInternal,
	})
}

// RunHubs runs the different broadcasting hubs
func (h *Handler) RunHubs() {
	// h.signalingHub.Run()
	// h.messagingHub.Run()
}

func getClaimes(c echo.Context) *jwt.Claims {
	user := c.Get("user").(*j.Token)
	return user.Claims.(*jwt.Claims)
}

func getToken(c echo.Context) *j.Token {
	return c.Get("user").(*j.Token)
}
