package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func createWebserver(port int, token string) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	
	e.POST("/receive", func (c echo.Context) error {
		if queryToken := c.QueryParam("token"); queryToken != token {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Empty or invalid token!",
			})
		}

		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Missing file from FormData",
			})
		}
		src, err := file.Open()
		if err != nil {
			return validateInternalError(c, err)
		}
		dst, err := os.Create(file.Filename)
		if err != nil {
			return validateInternalError(c, err)
		}
		defer src.Close()
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return validateInternalError(c, err)
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	})

	e.Logger.Fatal(e.Start("0.0.0.0:"+strconv.Itoa(port)))
	errGetGlobal <- errors.New("Webserver Closed")
}

func validateInternalError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"message": err.Error(),
	})
}
