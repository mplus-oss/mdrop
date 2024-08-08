package main

import (
	"flag"
	"fmt"
	"os"
)

const subCmdHelpMeesage string = `
Subcommand:
  auth
    subcommand
        Authenticate client to broker server
  get [options] <token>
    subcommand
        Create instance for retriving file from sender
  send [options] <file1> [file2] [file...]
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
	switch cmd {
	case "auth":
		AuthCommand()
	case "get":
		GetCommand(args)
	case "send":
		SendCommand(args)
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
