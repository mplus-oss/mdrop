package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: mdrop-tunnel <subcommand> <args>")
		os.Exit(1)
	}

	privateTokenByte, err := os.ReadFile("/.token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	privateToken := string(privateTokenByte)
	if privateToken == "" {
		fmt.Println("Private token not provided")
		os.Exit(1)
	}

	jwtSaltByte, err := os.ReadFile("/.cert")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jwtSalt := string(jwtSaltByte)
	if jwtSalt == "" {
		fmt.Println("JWT Private not provided")
		os.Exit(1)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "auth":
		AuthCommand(privateToken, args)
	default:
		fmt.Println("Subcommand: auth")
	}
}
