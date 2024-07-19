package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	internal "github.com/mplus-oss/mdrop/client/internal"
)

func GetCommand(args []string) {
	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		expired = flag.Int("expired", 3, "Expired timeout in hours")
		port = flag.Int("port", 0, "Specified port on broker, default is random. Range of port is about 10k - 59k.")
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
	fmt.Println(c.Token)
	fmt.Println(*expired)
	fmt.Println(*port)
}

func randRange(min, max int) int {
    return rand.IntN(max-min) + min
}
