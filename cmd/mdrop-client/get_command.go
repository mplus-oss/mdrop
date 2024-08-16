package main

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mplus-oss/mdrop/internal"
)

func GetCommand(args []string) {
	errChan := make(chan error, 1)

	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		help		= flag.Bool("help", false, "Print this message")
		fileNameOpt = flag.String("output", "", "Set file or directory output")
		localPort	= flag.Int("local-port", 6000, "Specified sender port remoted on local")
	)
	flag.Parse(args)

	token := flag.Arg(0)
	if *help || token == "" {
		fmt.Println("Command: mdrop get [options] <token>")
		flag.Usage()
		os.Exit(1)
	}

	// Parse token
	sender := TokenTransferJSON{}
	err := sender.ParseToken(token, &sender)
	if err != nil {
		internal.PrintErrorWithExit("getParseTokenError", err, 1)
	}

	// Parse Config File
	var authFile internal.ConfigSourceAuth
	err = authFile.ParseConfig(&authFile)
	if err != nil {
		internal.PrintErrorWithExit("getParseConfigError", err, 1)
	}
	if len(authFile.ListConfiguration) == 0 {
		internal.PrintErrorWithExit("getConfigFileEmpty", errors.New("Config file empty. Please log in using `mdrop auth` before executing this command."), 1)
	}
	config := authFile.ListConfiguration[authFile.Default]
	if sender.Host != config.Host {
		internal.PrintErrorWithExit("getHostNotMatch", errors.New("Host not match"), 1)
	}

	// Create tunnel before remote
	fmt.Println("Connecting to tunnel for fetch port...")
	sshConfig := internal.SSHParameter{
		ConfigFile: config,
		Command:    []string{"connect"},
		LocalPort:  *localPort,
		RemotePort: sender.RemotePort,
		IsRemote:   false,
	}
	go func() {
		_, err := StartShellTunnel(sshConfig, true)
		if err != nil {
			errChan <- err
		}
	}()

	// Check tunnel
	fmt.Println("Connecting to sender...")
	GetTcpReadyConnect(*localPort)

	// Downloading file
	isSingleFile := len(sender.Files) < 2
	for _, uuid := range sender.Files {
		// No error checking needed.
		fileName := GetPrompt(*localPort, uuid, *fileNameOpt, isSingleFile)
		checksum := GetChecksum(*localPort, uuid)
		filePath := GetDownload(*localPort, uuid, fileName, *fileNameOpt, isSingleFile)

		// Check checksum
		fmt.Println("Checking checksum...")
		fileDownloaded, err := os.Open(filePath)
		if err != nil {
			internal.PrintErrorWithExit("checksumFileOpen", err, 1)
		}
		hash := sha256.New()
		if _, err := io.Copy(hash, fileDownloaded); err != nil {
			internal.PrintErrorWithExit("checksumHashSum", err, 1)
		}
		checksumLocal := fmt.Sprintf("%x", hash.Sum(nil))
		if checksumLocal != checksum {
			internal.PrintErrorWithExit("checksumMismatch", errors.New("Checksum mismatch with sender"), 1)
		}
		fileDownloaded.Close()
	}

	fmt.Println("\nDownload success!")
}
