package main

import (
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
	flag := flag.NewFlagSet("mdrop get", flag.ExitOnError)
	var (
		help      = flag.Bool("help", false, "Print this message")
		output	  = flag.String("output", "", "Set output file name")
		localPort = flag.Int("localPort", 6000, "Specified sender port remoted on local")
	)
	flag.Parse(args)

	if *help {
		fmt.Println("Command: mdrop send [options] <file>")
		flag.Usage()
		os.Exit(1)
	}

	// TODO: Create tunnel before remote
	fmt.Println("Connecting to sender...")
	client := http.Client{}
	resp, err := client.Post(
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

	fileName := resp.Header.Get("X-Attachment-Name")
	if fileName == "" {
		internal.PrintErrorWithExit("sendHttpClientInvalidAttachmentName", err, 1)
	}
	if *output != "" {
		fileName = *output
	}
	fmt.Println("File found:", fileName)

	file, err := os.Create(fileName)
	if err != nil {
		internal.PrintErrorWithExit("sendFileCreation", err, 1)
	}

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
}
