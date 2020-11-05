// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUserServers returns all servers the user is a member of
func (h *Handler) GetUserServers(c echo.Context) error {
	claims := getClaimes(c)
	servers, err := h.userService.GetUserServers(claims.Username)
	if err != nil {
		log.Printf("ERROR [Router] Failed to get servers: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"servers": servers,
	})
}
