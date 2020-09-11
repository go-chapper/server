// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package handlers provides HTTP handlers
package handlers

import (
	"git.web-warrior.de/go-chapper/server/internal/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateAvatar(c echo.Context) error {
	return nil
}

func (h *Handler) GetAvatar(c echo.Context) error {
	p := utils.Join(h.config.Router.AvatarPath, c.Param("size"), c.Param("name")+".jpg")
	return c.File(p)
}
