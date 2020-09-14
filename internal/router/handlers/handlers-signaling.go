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

func (h *Handler) GetSignalingChannel(c echo.Context) error {
	conn, err := h.hub.CreateConnection(c.Response(), c.Request())
	if err != nil {
		log.Printf("ERROR [Router] Failed to upgrade connection: %v\n", err)
		return err
	}

	go h.hub.Register(conn)
	go conn.ListenRead()
	go conn.ListenWrite()

	return nil
}

func (h *Handler) GetSignalingToken(c echo.Context) error {
	claims := getClaimes(c)

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
