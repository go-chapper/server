// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"log"
	"net/http"

	"chapper.dev/server/internal/models"

	"github.com/labstack/echo/v4"
)

// DeleteDirect deletes a direct room identified by it's name
func (h *Handler) DeleteDirect(c echo.Context) error {
	return nil
}

// UpdateDirect updates a direct room
func (h *Handler) UpdateDirect(c echo.Context) error {
	return nil
}

// CreateDirect creates a direct room
func (h *Handler) CreateDirect(c echo.Context) error {
	claims := getClaimes(c)

	direct := new(models.Direct)
	err := c.Bind(direct)
	if err != nil {
		log.Printf("WARNING [Router] Unable to bind to model: %v\n", err)
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrBind,
		})
	}

	if direct.IsEmpty() {
		log.Println("WARNING [Router] Missing/empty data to create direct")
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrEmptyData,
		})
	}

	// TODO <2020/13/09>: Add check if callee is blocked and/or the caller is blocked by callee

	t, err := h.hub.Token(claims.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrCreateDirect,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"token": t,
	})
}

// GetDirect returns a direct room identified by it's name
func (h *Handler) GetDirect(c echo.Context) error {
	return nil
}
