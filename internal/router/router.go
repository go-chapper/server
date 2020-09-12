// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package router provides the top-level router
package router

import (
	"fmt"
	"log"

	"chapper.dev/server/internal/config"
	"chapper.dev/server/internal/modules/jwt"
	"chapper.dev/server/internal/router/handlers"
	"chapper.dev/server/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	l "github.com/labstack/gommon/log"
)

// Router is the top-level router instance wrapping it's dependencies
type Router struct {
	config *config.Config
	echo   *echo.Echo
}

// New creates a new router instance and returns it
func New(c *config.Config) *Router {
	e := echo.New()

	// Set debug mode (only for development)
	e.Debug = false

	// Set log level
	e.Logger.SetLevel(l.ERROR)

	// Hide startup message
	e.HideBanner = true

	// Register middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Enable GZIP compression
	if c.Router.EnableGZIP {
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}))
	}

	// Redirect www to non www
	e.Pre(middleware.NonWWWRedirect())

	return &Router{
		config: c,
		echo:   e,
	}
}

// AddRoutes adds all routes to the router instance and registers the handlers
func (r *Router) AddRoutes(handle *handlers.Handler) {
	// TODO: Move this to config validation / default values
	webRoot, err := utils.Abs(r.config.Router.WebPath)
	if err != nil {
		log.Panicln(err)
	}
	r.config.Router.WebPath = webRoot

	//// UNPROTECTED ROUTES ////
	// SPA
	r.echo.Use(middleware.Static(webRoot))

	// INVITE
	r.echo.GET("/i/:invite", handle.Invite)

	// MEDIA
	media := r.echo.Group("/media")
	media.GET("/images", handle.GetImage)
	media.GET("/videos", handle.GetVideo)

	// JWT middleware setup
	jwtware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(r.config.Router.JWTSigningKey),
		Claims:     &jwt.Claims{},
	})

	// AVATAR
	avatar := r.echo.Group("/avatar")
	avatar.POST("/:name", handle.UpdateAvatar, jwtware)
	avatar.GET("/:size/:name", handle.GetAvatar)

	// PUBLIC KEY
	key := r.echo.Group("/key")
	key.GET("/:username", handle.GetKey, jwtware)

	//// API ////
	api := r.echo.Group("/api", jwtware)
	v1 := api.Group("/v1")

	// INVITES
	invite := v1.Group("/invite")
	invite.DELETE("/:name", handle.DeleteInvite)
	invite.PUT("", handle.CreateInvite)

	// PROFILE
	profile := v1.Group("/profile")
	profile.GET("/:username", handle.GetProfile)

	// VIRTUAL SERVERS
	server := v1.Group("/servers")
	server.DELETE("/:server-hash", handle.DeleteServer)
	server.POST("/:server-hash", handle.UpdateServer)
	server.GET("/:server-hash", handle.GetServer)
	server.PUT("", handle.CreateServer)
	server.GET("", handle.GetServers)

	// ROOMS
	server.DELETE("/:server-hash/rooms/:room-hash", handle.DeleteRoom)
	server.POST("/:server-hash/rooms/:room-hash", handle.UpdateRoom)
	server.GET("/:server-hash/rooms/:room-hash", handle.GetRoom)
	server.PUT("/:server-hash/rooms", handle.CreateRoom)
	server.GET("/:server-hash/rooms", handle.GetRooms)

	// DIRECT
	direct := v1.Group("/direct")
	direct.DELETE("/:direct-hash", handle.DeleteDirect)
	direct.POST("/:direct-hash", handle.UpdateDirect)
	direct.GET("/:direct-hash", handle.GetDirect)
	direct.PUT("", handle.CreateDirect)

	// AUTH
	auth := r.echo.Group("/auth")
	auth.POST("/code/register", handle.AuthRegisterCode)
	auth.POST("/register", handle.AuthRegister)
	auth.POST("/refresh", handle.AuthRefresh)
	auth.POST("/login", handle.AuthLogin)
	auth.POST("/code", handle.AuthCode)

	// This serves the correct SPA route (even when reloading)
	r.echo.File("/*", webRoot)
}

// Run starts the HTTP Server or returns an error
func (r *Router) Run() error {
	return r.echo.Start(fmt.Sprintf(":%d", r.config.Router.Port))
}
