// Copyright (c) 2020-present Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, h.config.Router.Domain)
}
