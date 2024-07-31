package main

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func SendWebserver(localPort int, file string) error {
    e := echo.New()
    e.HidePort = true
    e.HideBanner = true

    e.GET("/receive", func(c echo.Context) (err error) {
        err = c.File(file)
        return err
    })

    return e.Start("0.0.0.0:"+strconv.Itoa(localPort))
}
