package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mplus-oss/mdrop/internal"
)

func SendCommand(args []string) {
	var err error
	var errChan chan error = make(chan error)

	flag := flag.NewFlagSet("mdrop send", flag.ExitOnError)
	var (
		help      = flag.Bool("help", false, "Print this message")
		localPort = flag.Int("localPort", 6000, "Specified sender port on local")
	)
	flag.Parse(args)

	file := flag.Arg(0)
	if *help || file == "" {
		fmt.Println("Command: mdrop send [options] <file>")
		flag.Usage()
		os.Exit(1)
	}

	// Deploy webserver instance
	go func() {
		fmt.Println("Spawning webserver...")
		err := SendWebserver(*localPort, file)
		if err != nil {
			errChan <- err
		}
	}()

	// Check if there's some error on different threads
	err = <-errChan
	if err != nil {
		internal.PrintErrorWithExit("receiveWebserverFatalError", err, 1)
	}
}
