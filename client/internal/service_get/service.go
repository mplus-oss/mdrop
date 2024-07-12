package service_get

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func StartGetService() {
	// Server stuff
	e := echo.New()
	e.HideBanner = true

	// Routing
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})
	// x-url-form-encoded
	e.POST("/transfer", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		homedir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		time := time.Now()
		dst, err := os.Create(fmt.Sprintf(
			"%v/Downloads/%v-%v",
			homedir,
			time.UTC().UnixMilli(),
			file.Filename,
		))
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.Logger.Fatal(e.Start("0.0.0.0:7000"))
}
