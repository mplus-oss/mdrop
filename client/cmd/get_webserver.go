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
	"github.com/schollz/progressbar/v3"
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
		fmt.Print("Accept? [(Y)es/(N)o] -> ")
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

		// Validate if the data can be override or not
		currentPath, err := os.Getwd()
		if err != nil {
			return validateInternalError(c, err)
		}
		currentPath += "/"+file.Filename
		if fileStatus, _ := os.Stat(currentPath); fileStatus != nil {
			fmt.Print("There's duplicate file. Action? [(Y)es/(N)o/(R)ename] -> ")
			prompt, err = reader.ReadString('\n')
			if err != nil {
				return validateInternalError(c, err)
			}
			prompt = strings.Replace(prompt, "\n", "", -1)
			if strings.ToLower(prompt) == "n" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "File declined from receiver",
				})
			}
			if strings.ToLower(prompt) == "r" {
				fmt.Print("Change filename ["+file.Filename+"] -> ")
				prompt, err = reader.ReadString('\n')
				if err != nil {
					return validateInternalError(c, err)
				}
				prompt = strings.Replace(prompt, "\n", "", -1)
				if prompt == file.Filename {
					fmt.Println("Canceled. Duplicate file.")
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"message": "Duplicate file from receiver",
					})
				}
				file.Filename = prompt
			}
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
		bar := progressbar.DefaultBytes(file.Size, file.Filename)
		if _, err = io.Copy(io.MultiWriter(dst, bar), src); err != nil {
			return validateInternalError(c, err)
		}

		// Complete
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
