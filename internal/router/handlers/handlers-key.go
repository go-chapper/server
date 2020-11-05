// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetKey returns the public key of a given user
func (h *Handler) GetKey(c echo.Context) error {
	key, err := h.userService.GetUserPublicKey(c.Param("username"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, Map{
				"error": ErrGetUser,
			})
		}

		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.String(http.StatusOK, key)
}
