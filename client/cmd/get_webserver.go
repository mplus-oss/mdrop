package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func createWebserver(port int) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	
	e.GET("/receive", func (c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	})

	e.Logger.Fatal(e.Start("0.0.0.0:"+strconv.Itoa(port)))
	errGetGlobal <- errors.New("Webserver Closed")
}
