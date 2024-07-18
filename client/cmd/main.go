package main

import (
	"flag"
	"fmt"
	"os"
)

const helpMeesage string = `Command: mdrop <command> [options]

Usage of mdrop:
  auth [--url=uri] [--token=random_string]
    subcommand
        Authenticate client to broker server
  get
    subcommand
        Create instance for retriving file from sender
  send
    subcommand
        Send file to reciever instance`

func main() {
	help := flag.Bool("help", false, "Print this message")
	flag.Parse()

	args := flag.Args()
	if len(os.Args) == 1 {
		fmt.Println(helpMeesage)
		os.Exit(1)
	}

	if *help {
		fmt.Println(helpMeesage)
		os.Exit(1)
	}

	cmd, args := args[0], args[1:]
	switch(cmd) {
	case "auth":
		fmt.Println("Not implemented")
		os.Exit(1)
	case "get":
		fmt.Println("Not implemented")
		os.Exit(1)
	case "send":
		fmt.Println("Not implemented")
		os.Exit(1)
	default:
		fmt.Println(helpMeesage)
		os.Exit(1)
	}
}

// For references
func subCommand(args []string) {
	flag := flag.NewFlagSet("mdrop sb", flag.ExitOnError)
	var (
		url = flag.String("url", "https://example.com", "URL of broker")
	)
	flag.Parse(args)

	fmt.Println(*url)
}
