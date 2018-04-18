package application

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	static "github.com/Code-Hex/echo-static"
)

func (a *Application) setupServer() {
	a.Server.Pre(middleware.RemoveTrailingSlash())
	a.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	a.Server.Use(middleware.Recover())
	a.Server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "02-01-2006 | 03:04:05.00000",
		Format:           `[CFFT] ${time_custom} | ${method} | ${status} | ${uri} -> ${latency_human}` + "\n",
	}))
	a.Server.HideBanner = true
	a.Server.HidePort = true
	a.Logger.Debug("Web Server middleware setup complete.")
}

func (a *Application) setupErrorHandler() {
	a.Server.HTTPErrorHandler = a.customHTTPErrorHandler
	a.Logger.Debug("Web Server error handler setup complete.")
}

func (a *Application) setupRoutes() {
	a.Server.GET("/ping", a.routesHomePing)
	a.Server.GET("/health", func(e echo.Context) error {
		a.Checkers.Handler.ServeHTTP(e.Response().Writer,e.Request())
		return nil
	})
	a.Server.GET("/", a.routesHomeIndex)
	a.Logger.Debug("Server core routes initialisation complete.")
}

func (a *Application) setupSetupRoutes() {
	a.Server.GET("/setup", a.routesHomeSetup1)
	a.Server.GET("/hxsetup1", a.routesHomeHXSetup1)
	a.Server.POST("/hxsetup2", a.routesHomeHXSetup2)
	a.Server.POST("/hxsetup3", a.routesHomeHXSetup3)
	a.Server.POST("/hxsetup4", a.routesHomeHXSetup4)
	a.Server.POST("/hxsetup5", a.routesHomeHXSetup5)
	a.Logger.Debug("Server setup routes initialisation complete.")
}

func (a *Application) setupTemplates() {
	a.Server.Use(static.ServeRoot("/", NewAssets("public")))
	a.Server.Renderer = NewTemplate()
	a.Logger.Debug("Server Renderer initialised successfully.")
}