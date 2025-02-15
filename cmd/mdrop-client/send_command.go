package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/google/uuid"
	"github.com/mplus-oss/mdrop/internal"
)

func SendCommand(args []string) {
	var err error
	errChan := make(chan error, 0)
	outputChan := make(chan string, 0)

	flag := flag.NewFlagSet("mdrop send", flag.ExitOnError)
	var (
		help      = flag.Bool("help", false, "Print this message")
		localPort = flag.Int("local-port", 6000, "Specified sender port on local")
	)
	flag.Parse(args)

	files := flag.Args()
	if *help || len(files) == 0 {
		fmt.Println("Command: mdrop send [options] <file1> [file2] [file...]")
		flag.Usage()
		os.Exit(1)
	}
	// Parsing glob
	files, err = ParseGlob(files)
	if err != nil {
		internal.PrintErrorWithExit("sendParseGlobError", err, 1)
	}

	// Parse Config
	var authFile internal.ConfigSourceAuth
	err = authFile.ParseConfig(&authFile)
	if err != nil {
		internal.PrintErrorWithExit("sendParseConfigError", err, 1)
	}
	if len(authFile.ListConfiguration) == 0 {
		internal.PrintErrorWithExit("getConfigFileEmpty", errors.New("Config file empty. Please log in using `mdrop auth` before executing this command."), 1)
	}
	config := authFile.ListConfiguration[authFile.Default]

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

	// Deploy webserver instance & add uuid for resolver
	go func() {
		fileUUID := []string{}

		for range files {
			fileUUID = append(fileUUID, uuid.New().String())
		}

		// Print token
		token, err := TokenTransferJSON{
			Host:       config.Host,
			RemotePort: remotePort,
			Files:      fileUUID,
		}.GenerateToken()
		if err != nil {
			errChan <- err
		}

		fmt.Print("\nPlease copy this token to the receiver.")
		fmt.Print("\nToken: " + token + "\n\n")

		fmt.Println("Token copied to clipboard!")
		err = clipboard.WriteAll(token)
		if err != nil {
			fmt.Println("Warning: Failed to copy token to clipboard.")
		}

		fmt.Println("Spawning webserver...")
		err = SendWebserver(*localPort, files, fileUUID)
		errChan <- err
	}()

	// Check if there's some error on different threads
	err = <-errChan
	if err != nil {
		internal.PrintErrorWithExit("receiveWebserverFatalError", err, 1)
	}
	fmt.Println("Done!")
}
