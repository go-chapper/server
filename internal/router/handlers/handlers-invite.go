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

// Invite handles the invite when a user enters an invite link
func (h *Handler) Invite(c echo.Context) error {
	return c.String(http.StatusOK, c.Param("invite"))
}

// CreateInvite creates an invite link and returns it
func (h *Handler) CreateInvite(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanCreateInvite {
		return c.JSON(http.StatusUnauthorized, Map{
			"errror": "Invalid request",
			"code":   ErrUnauthorized,
		})
	}

	// TODO <2020/08/09>: Dont hardcode server name, get it from request
	invite, err := h.inviteService.CreateInvite(claims.Username, "Chapper", false)
	if err != nil {
		log.Printf("ERROR [Router] Failed to create invite: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Map{
			"errror": "Internal server error",
			"code":   ErrCreateInvite,
		})
	}

	return c.JSON(http.StatusOK, Map{
		"inviteUrl": invite.URLString,
	})
}

// DeleteInvite deletes an invite link indentified by it's name
func (h *Handler) DeleteInvite(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanDeleteInvite {
		return c.JSON(http.StatusUnauthorized, Map{
			"errror": "Invalid request",
			"code":   ErrUnauthorized,
		})
	}

	return nil
}
