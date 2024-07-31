package main

import (
	"flag"
	"fmt"
	"os"
)

const subCmdHelpMeesage string = `
Subcommand:
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
		printUsage()
	}

	if *help {
		printUsage()
	}

	cmd, args := args[0], args[1:]
	switch(cmd) {
	case "auth":
		fmt.Println("Not implemented")
	case "get":
		fmt.Println("Not implemented")
	case "send":
        SendCommand(args)
		fmt.Println("Not implemented")
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Print("Command: mdrop <command> [options]\n\n")
	flag.Usage()
	fmt.Println(subCmdHelpMeesage)
	os.Exit(1)
}
