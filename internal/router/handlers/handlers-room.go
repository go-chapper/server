// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"log"
	"net/http"
	"strings"

	"git.web-warrior.de/go-chapper/server/internal/models"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

// CreateRoom creates a room
// NOTE(Techassi): This silently fails if we try to create a room with the same name/hash
// but doesn't create a new one. Sooooo... Success I guess
func (h *Handler) CreateRoom(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanCreateRoom {
		return c.JSON(http.StatusUnauthorized, Map{
			"errror": "Invalid request",
			"code":   ErrUnauthorized,
		})
	}

	room := new(models.Room)
	err := c.Bind(room)
	if err != nil {
		log.Printf("WARNING [Router] Unable to bind to model: %v\n", err)
		return c.JSON(http.StatusBadRequest, Map{
			"error": "Invalid request",
			"code":  ErrBind,
		})
	}

	if room.IsEmpty() {
		log.Println("WARNING [Router] Missing/empty data to create room")
		return c.JSON(http.StatusBadRequest, Map{
			"error": "Invalid request",
			"code":  ErrEmptyData,
		})
	}

	if room.Invalid() {
		log.Println("WARNING [Router] Invalid data to create room")
		return c.JSON(http.StatusBadRequest, Map{
			"error": "Invalid request",
			"code":  ErrInvalidData,
		})
	}

	err = h.roomService.CreateRoom(c.Param("server-hash"), room)
	if err != nil {
		log.Printf("ERROR [Router] Failed to create room: %v\n", err)

		// TODO <2020/10/09>: Optimize this FOR SURE
		if strings.HasPrefix(err.Error(), "Error 1062") {
			return c.JSON(http.StatusBadRequest, Map{
				"error": "Invalid request",
				"code":  ErrRoomnameTaken,
			})
		}

		return c.JSON(http.StatusInternalServerError, Map{
			"errror": "Internal server error",
			"code":   ErrCreateRoom,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": "Success",
		"code":   StatusRoomCreated,
	})
}

// GetRoom returns a room identified by it's name
func (h *Handler) GetRoom(c echo.Context) error {
	room, err := h.roomService.GetRoom(c.Param("server-hash"), c.Param("room-hash"))
	if err != nil {
		log.Printf("ERROR [Router] Failed to get server: %v\n", err)

		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, Map{
				"error": "Invalid request",
				"code":  "",
			})
		}

		return c.JSON(http.StatusInternalServerError, Map{
			"errror": "Internal server error",
			"code":   "",
		})
	}

	return c.JSON(http.StatusOK, Map{
		"room": room,
	})
}

// GetRooms returns a room identified by it's name
func (h *Handler) GetRooms(c echo.Context) error {
	rooms, err := h.roomService.GetRooms(c.Param("server-hash"))
	if err != nil {
		log.Printf("ERROR [Router] Failed to get rooms: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"errror": "Internal server error",
			"code":   "",
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
	return nil
}
