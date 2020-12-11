// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"net/http"

	"chapper.dev/server/internal/services"

	"github.com/labstack/echo/v4"
)

// AuthRegister registers a new user
func (h *Handler) AuthRegister(c echo.Context) error {
	err := h.authService.Register(c)
	if err != nil {
		if se, ok := err.(*services.ServiceError); ok {
			h.logger.Errorf("Router: %v", se)
			return c.JSON(se.Code(), Map{
				"error": se.Err(),
			})
		}

		h.logger.Errorf("Router: %v", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": StatusRegistered,
	})
}

// AuthRefresh refreshes the JWT token
func (h *Handler) AuthRefresh(c echo.Context) error {
	return nil
}

// AuthLogin logs a user in
func (h *Handler) AuthLogin(c echo.Context) error {
	token, err := h.authService.Login(c)
	if err != nil {
		if se, ok := err.(*services.ServiceError); ok {
			h.logger.Errorf("Router: %v", se)
			return c.JSON(se.Code(), Map{
				"error": se.Err(),
			})
		}

		h.logger.Errorf("Router: %v", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	// If the token is empty and there is no error we need to validate the 2FA code
	if token == "" {
		return c.JSON(http.StatusOK, Map{
			"state": "code",
			"token": token,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"state": "authenticated",
		"token": token,
	})
}
