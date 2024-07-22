package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func createWebserver(port int, token string) {
	reader := bufio.NewReader(os.Stdin)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	
	e.POST("/receive", func (c echo.Context) error {
		// Get query token
		if queryToken := c.QueryParam("token"); queryToken != token {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Empty or invalid token!",
			})
		}

		// Get file
		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Missing file from FormData",
			})
		}

		// Prompting
		fmt.Println("\nIncoming file:", file.Filename)
		fmt.Print("Accept? [Y/N] -> ")
		prompt, err := reader.ReadString('\n')
		if err != nil {
			return validateInternalError(c, err)
		}
		prompt = strings.Replace(prompt, "\n", "", -1)
		if strings.ToLower(prompt) != "y" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "File declined from receiver",
			})
		}

		// File scanning
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

		// Copy file
		pr := &Progress{
			TotalSize: file.Size,
		}
		if _, err = io.Copy(dst, io.TeeReader(src, pr)); err != nil {
			return validateInternalError(c, err)
		}

		// Complete
		fmt.Println(file.Filename, "successfully transferred!")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	})

	e.Logger.Fatal(e.Start("0.0.0.0:"+strconv.Itoa(port)))
	errGetGlobal <- errors.New("Webserver Closed")
}

type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Aggregating length of upload progress on io.Writer instance.
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}

	fmt.Printf("File upload in progress: %d\n", pr.BytesRead)
}

func validateInternalError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"message": err.Error(),
	})
}
