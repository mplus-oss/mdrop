package main

import (
	"flag"
	"fmt"
	"os"
)

func GetCommand(args []string) {
	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		help = *flag.Bool("help", false, "Print this message")
	)
	flag.Parse(args)

	if help {
		fmt.Println("Command: mdrop send [options] <file>")
		flag.Usage()
		os.Exit(1)
	}
}
