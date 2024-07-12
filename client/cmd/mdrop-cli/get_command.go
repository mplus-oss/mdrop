package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func StartGetCommand() {
	// Server stuff
	e := echo.New()
	e.HideBanner = true

	// Routing
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	e.Logger.Fatal(e.Start("0.0.0.0:7000"))
}
