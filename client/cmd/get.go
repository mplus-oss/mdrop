package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	internal "github.com/mplus-oss/mdrop/client/internal"
)

var errGlobal chan error = make(chan error)

func GetCommand(args []string) {
	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		expired	  = flag.Int("expired", 3, "Expired timeout in hours.")
		port	  = flag.Int("port", 0, "Specified port on broker. Range of port is about 10k - 59k. (default rand(10000, 59999))")
		localPort = flag.Int("localPost", 6000, "Specified port on local.")
	)
	flag.Parse(args)

	if *port == 0 {
		*port = randRange(10000, 59999)
	}
	if *port < 10000 || *port > 59999 {
		fmt.Println("Invalid port range!")
		os.Exit(1)
	}

	var c internal.ConfigFile
	err := c.ParseConfig(&c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print("Creating the room...\n\n")

	// Start the Echo server
	go createWebserver(*localPort)

	// Get token data
	var path = fmt.Sprintf(
		"%v/room/create?durationInHours=%v&port=%v",
		c.URL,
		*expired,
		*port,
	)
	roomData, err := getToken(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(roomData.Token)

	err = <- errGlobal
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
	errGlobal <- errors.New("Webserver Closed")
}

func getToken(path string) (CreateRoomJSONReturn, error) {
	var client = &http.Client{}
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return CreateRoomJSONReturn{} ,err
	}
	res, err := client.Do(req)
	if err != nil {
		return CreateRoomJSONReturn{} ,err
	}
	defer res.Body.Close()

	var roomData CreateRoomJSONReturn
	err = json.NewDecoder(res.Body).Decode(&roomData)
	if err != nil {
		return CreateRoomJSONReturn{} ,err
	}

	return roomData, nil
}

type CreateRoomJSONReturn struct {
	Token string `json:"token"`

	// This respon fired when the API is failed
	ErrorTitle string `json:"title"`
	Errors     struct {
		Port []string `json:"port"`
	} `json:"errors"`
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
