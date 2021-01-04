// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"

	"chapper.dev/server/internal/services/errors"

	"github.com/labstack/echo/v4"
)

// GetInvite handles the invite when a user enters an invite link
func (h *Handler) GetInvite(c echo.Context) error {
	return c.String(http.StatusOK, c.Param("invite"))
}

// CreateInvite creates an invite link and returns it
func (h *Handler) CreateInvite(c echo.Context) error {
	claims := getClaimes(c)

	if !claims.Privileges.CanCreateInvite {
		return c.JSON(http.StatusUnauthorized, Map{
			"error": ErrUnauthorized,
		})
	}

	invite, err := h.inviteService.CreateInvite(claims.Username, c)
	if err != nil {
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

	return c.JSON(http.StatusOK, Map{
		"url": invite.ToURL(h.config.Router.Domain),
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
