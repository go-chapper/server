// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"chapper.dev/server/internal/services/errors"

	"github.com/labstack/echo/v4"
)

// CreateRoom handles incoming requests to create a room
func (h *Handler) CreateRoom(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanCreateRoom {
		return c.JSON(http.StatusUnauthorized, Map{
			"error": ErrUnauthorized,
		})
	}

	err := h.roomService.CreateRoom(c)
	if err != nil {
		if se, ok := err.(*errors.ServiceError); ok {
			h.logger.Errorc(routerCtx, se)
			return c.JSON(se.Code(), Map{
				"error": se.Err(),
			})
		}

		h.logger.Errorc(routerCtx, err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": "created",
	})
}

// GetRoom handles incoming requests to get one room
func (h *Handler) GetRoom(c echo.Context) error {
	room, err := h.roomService.GetRoom(c)
	if err != nil {
		if se, ok := err.(*errors.ServiceError); ok {
			h.logger.Errorc(routerCtx, se)
			return c.JSON(se.Code(), Map{
				"error": se.Err(),
			})
		}

		h.logger.Errorc(routerCtx, err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"room": room,
	})
}

// GetRooms handles incoming requests to get multiple rooms
func (h *Handler) GetRooms(c echo.Context) error {
	rooms, err := h.roomService.GetRooms()
	if err != nil {
		if se, ok := err.(*errors.ServiceError); ok {
			h.logger.Errorc(routerCtx, se)
			return c.JSON(se.Code(), Map{
				"error": se.Err(),
			})
		}

		h.logger.Errorc(routerCtx, err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"rooms": rooms,
	})
}

// UpdateRoom handles incoming requests to update one room
func (h *Handler) UpdateRoom(c echo.Context) error {
	err := h.roomService.UpdateRoom(c)
	if err != nil {
		if se, ok := err.(*errors.ServiceError); ok {
			h.logger.Errorc(routerCtx, se)
			return c.JSON(se.Code(), Map{
				"error": se.Err(),
			})
		}

		h.logger.Errorc(routerCtx, err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": "updated",
	})
}

// DeleteRoom handles incoming requests to delete one room
func (h *Handler) DeleteRoom(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanDeleteRoom {
		return c.JSON(http.StatusUnauthorized, Map{
			"errror": ErrUnauthorized,
		})
	}

	err := h.roomService.DeleteRoom(c)
	if err != nil {
		if se, ok := err.(*errors.ServiceError); ok {
			h.logger.Errorc(routerCtx, se)
			return c.JSON(se.Code(), Map{
				"error": se.Err(),
			})
		}

		h.logger.Errorc(routerCtx, err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": "deleted",
	})
}
