// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"chapper.dev/server/internal/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateAvatar(c echo.Context) error {
	return nil
}

func (h *Handler) GetAvatar(c echo.Context) error {
	p := utils.Join(h.config.Router.AvatarPath, c.Param("size"), c.Param("name")+".jpg")
	return c.File(p)
}
