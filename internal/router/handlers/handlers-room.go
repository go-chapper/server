// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"log"
	"net/http"

	"chapper.dev/server/internal/services/errors"

	"github.com/labstack/echo/v4"
)

// CreateRoom creates a room
// NOTE(Techassi): This silently fails if we try to create a room with the same name/hash
// but doesn't create a new one. Sooooo... Success I guess
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
		"status": StatusRoomCreated,
	})
}

// GetRoom returns a room identified by it's name
func (h *Handler) GetRoom(c echo.Context) error {
	room, err := h.roomService.GetRoom(c.Param("room-hash"))
	if err != nil {
		log.Printf("ERROR [Router] Failed to get server: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"errror": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"room": room,
	})
}

// GetRooms returns a room identified by it's name
func (h *Handler) GetRooms(c echo.Context) error {
	rooms, err := h.roomService.GetRooms()
	if err != nil {
		log.Printf("ERROR [Router] Failed to get rooms: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"errror": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"rooms": rooms,
	})
}

// UpdateRoom updates a room
func (h *Handler) UpdateRoom(c echo.Context) error {
	return nil
}

// DeleteRoom deletes a room identified by it's name
func (h *Handler) DeleteRoom(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanDeleteRoom {
		return c.JSON(http.StatusUnauthorized, Map{
			"errror": ErrUnauthorized,
		})
	}

	err := h.roomService.DeleteRoom(c.Param("room-hash"))
	if err != nil {
		log.Printf("ERROR [Router] Failed to delete room: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"errror": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": StatusRoomDeleted,
	})
}
