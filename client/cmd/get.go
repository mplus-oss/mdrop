package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"

	internal "github.com/mplus-oss/mdrop/client/internal"
)

func GetCommand(args []string) {
	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		expired = flag.Int("expired", 3, "Expired timeout in hours.")
		port    = flag.Int("port", 0, "Specified port on broker. Range of port is about 10k - 59k. (default rand(10000, 59999))")
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

	var client = &http.Client{}
	fmt.Println("Creating the room...")
	var pathStr = c.URL+"/room/create?durationInHours="+strconv.Itoa(*expired)+"&port="+strconv.Itoa(*port)
	req, err := http.NewRequest("POST", pathStr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	var roomData CreateRoomJSONReturn
	err = json.NewDecoder(res.Body).Decode(&roomData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(roomData.Token)
}

type CreateRoomJSONReturn struct {
	Token string `json:"token"`

	// This respon fired when the API is failed
	ErrorTitle string `json:"title"`
	Errors     struct {
		Port []string `json:"port"`
	} `json:"errors"`
}

func gatherToken() string {
	return ""
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
