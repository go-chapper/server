// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"log"
	"net/http"
	"strings"

	"chapper.dev/server/internal/models"

	"github.com/labstack/echo/v4"
)

// CreateServer creates a server
func (h *Handler) CreateServer(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanCreateServer {
		return c.JSON(http.StatusUnauthorized, Map{
			"error": ErrUnauthorized,
		})
	}

	server := new(models.Server)
	err := c.Bind(server)
	if err != nil {
		log.Printf("WARNING [Router] Unable to bind to model: %v\n", err)
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrBind,
		})
	}

	if server.IsEmpty() {
		log.Println("WARNING [Router] Missing/empty data to create server")
		return c.JSON(http.StatusBadRequest, Map{
			"error": ErrEmptyData,
		})
	}

	err = h.serverService.CreateServer(server)
	if err != nil {
		log.Printf("ERROR [Router] Failed to create server: %v\n", err)

		// TODO <2020/10/09>: Optimize this FOR SURE
		if strings.HasPrefix(err.Error(), "Error 1062") {
			return c.JSON(http.StatusBadRequest, Map{
				"error": ErrServernameTaken,
			})
		}

		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrCreateServer,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": StatusServerCreated,
	})
}

// GetServer returns a server identified by it's hash
func (h *Handler) GetServer(c echo.Context) error {
	server, err := h.serverService.GetServer(c.Param("server-hash"))
	if err != nil {
		log.Printf("ERROR [Router] Failed to get server: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"server": server,
	})
}

// GetServers returns all servers
func (h *Handler) GetServers(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanSeeAllServers {
		return c.JSON(http.StatusUnauthorized, Map{
			"error": ErrUnauthorized,
		})
	}

	servers, err := h.serverService.GetServers()
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

// UpdateServer updates a server
func (h *Handler) UpdateServer(c echo.Context) error {
	return nil
}

// DeleteServer deletes a server identified by it's hash
func (h *Handler) DeleteServer(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanDeleteServer {
		return c.JSON(http.StatusUnauthorized, Map{
			"error": ErrUnauthorized,
		})
	}

	err := h.serverService.DeleteServer(c.Param("server-hash"))
	if err != nil {
		log.Printf("ERROR [Router] Failed to delete server: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"error": ErrInternal,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"status": "Success",
		"code":   "",
	})
}
