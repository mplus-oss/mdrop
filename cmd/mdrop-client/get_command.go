package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mplus-oss/mdrop/internal"
	"github.com/schollz/progressbar/v3"
)

func GetCommand(args []string) {
	errChan := make(chan error, 1)

	reader := bufio.NewReader(os.Stdin)
	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		help		= flag.Bool("help", false, "Print this message")
		fileNameOpt = flag.String("file", "", "Set filename")
		localPort	= flag.Int("localPort", 6000, "Specified sender port remoted on local")
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
	var config internal.ConfigFile
	err = config.ParseConfig(&config)
	if err != nil {
		internal.PrintErrorWithExit("getParseConfigError", err, 1)
	}
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
	client := http.Client{}
	resp, err := client.Get(
		fmt.Sprintf("http://localhost:%v/checksum", *localPort),
	)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClientChecksum", err, 1)
	}
	if resp.StatusCode != http.StatusOK {
		internal.PrintErrorWithExit("sendHttpClientResponseChecksum", err, 1)
	}
	checksumBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClientReadChecksum", err, 1)
	}
	checksum := string(checksumBytes)

	resp, err = client.Post(
		fmt.Sprintf("http://localhost:%v/receive", *localPort),
		"binary/octet-stream",
		nil,
	)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClient", err, 1)
	}
	if resp.StatusCode != http.StatusOK {
		internal.PrintErrorWithExit("sendHttpClientResponse", err, 1)
	}
	defer resp.Body.Close()

	// Set filename from header or from output
	fileName := resp.Header.Get("X-Attachment-Name")
	if fileName == "" {
		internal.PrintErrorWithExit("sendHttpClientInvalidAttachmentName", err, 1)
	}
	if *fileNameOpt != "" {
		fileName = *fileNameOpt
	}
	fmt.Println("File found:", fileName)

	// Check if there's duplicate file
	filePath, err := os.Getwd()
	if err != nil {
		internal.PrintErrorWithExit("sendFileWorkDir", err, 1)
	}
	if fileStatus, _ := os.Stat(filePath+"/"+fileName); fileStatus != nil {
		fmt.Print("There's duplicate file. Action? [(R)eplace/R(e)name/(C)ancel] [Default: R] -> ")
		prompt, err := reader.ReadString('\n')
		if err != nil {
			internal.PrintErrorWithExit("sendPromptError", err, 1)
		}
		prompt = strings.Replace(prompt, "\n", "", -1)
		if strings.ToLower(prompt) == "c" {
			internal.PrintErrorWithExit("sendPromptCancel", errors.New("Canceled by action"), 0)
		}
		if strings.ToLower(prompt) == "e" {
			fmt.Print("Change filename ["+fileName+"]: ")
			prompt, err = reader.ReadString('\n')
			if err != nil {
				internal.PrintErrorWithExit("sendPromptError", err, 1)
			}
			prompt = strings.Replace(prompt, "\n", "", -1)
			if prompt == fileName {
				internal.PrintErrorWithExit("sendPromptDuplicateFilename", errors.New("Canceled by action"), 0)
			}
			fileName = prompt
		}
	}
	filePath += "/"+fileName

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		internal.PrintErrorWithExit("sendFileCreation", err, 1)
	}

	// Downloading file
	fmt.Println("Downloading...")
	bar := progressbar.DefaultBytes(resp.ContentLength, fileName)
	_, err = io.Copy(io.MultiWriter(bar, file), resp.Body)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EOF") {
			err = errors.New("Broken pipe from sender because forced close or terminated.")
		}
		internal.PrintErrorWithExit("sendStreamFile", err, 1)
	}

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

	fmt.Println("Download success!")
}
