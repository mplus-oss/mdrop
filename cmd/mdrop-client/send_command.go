package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mplus-oss/mdrop/internal"
)

func SendCommand(args []string) {
	var err error
	errChan := make(chan error, 0)
	outputChan := make(chan string, 0)

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

	// Parse Config
	var config internal.ConfigFile
	err = config.ParseConfig(&config)
	if err != nil {
		internal.PrintErrorWithExit("sendParseConfigError", err, 1)
	}

	// Get Port from Tunnel
	fmt.Println("Connecting to tunnel for fetch port...")
	sshConfig := internal.SSHParameter{
		ConfigFile: config,
		Command:    []string{"get-port"},
	}
	go func() {
		sshOutput, err := StartShellTunnel(sshConfig, false)
		if err != nil {
			outputChan <- "Nil"
			errChan <- err
		}
		if strings.Contains(sshOutput, "Get-Port-Found") {
			sshOutput = strings.Split(sshOutput, " ")[1]
			outputChan <- sshOutput
			errChan <- nil
		} else {
			outputChan <- "Nil"
			errChan <- errors.New("Unexpected output from Tunnel")
		}
	}()
	remotePort, err := strconv.Atoi(<-outputChan)
	err = <-errChan
	if err != nil {
		internal.PrintErrorWithExit("authTunnelError", err, 1)
	}

	// Deploy webserver tunnel
	sshConfig.IsRemote = true
	sshConfig.RemotePort = remotePort
	sshConfig.LocalPort = *localPort
	sshConfig.Command = []string{"connect"}
	go func() {
		_, err := StartShellTunnel(sshConfig, true)
		if err != nil {
			errChan <- err
		}
	}()

	// Deploy webserver instance
	go func() {
		fmt.Println("Spawning webserver...")
		err := SendWebserver(*localPort, file)
		errChan <- err
	}()

	// Check if there's some error on different threads
	err = <-errChan
	if err != nil {
		internal.PrintErrorWithExit("receiveWebserverFatalError", err, 1)
	}
	fmt.Println("Done!")
}
