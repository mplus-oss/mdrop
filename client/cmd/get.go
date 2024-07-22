package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	internal "github.com/mplus-oss/mdrop/client/internal"
)

var errGetGlobal chan error = make(chan error)

func GetCommand(args []string) {
	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		expired	  = flag.Int("expired", 3, "Expired timeout in hours.")
		port			= flag.Int("port", 0, "Specified port on broker. Range of port is about 10k - 59k. (default rand(10000, 59999))")
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
	fmt.Println(
		"\nCopy this token to the sender. Make sure sender has authenticated to",
		c.URL,
		"before sending the file.",
	)

	err = <- errGetGlobal
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
