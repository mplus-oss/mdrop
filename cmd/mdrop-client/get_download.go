package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mplus-oss/mdrop/internal"
	"github.com/schollz/progressbar/v3"
)

func GetChecksum(localPort int, uuid string) string {
	client := http.Client{}
	resp, err := client.Get(
		fmt.Sprintf("http://localhost:%v/checksum-%v", localPort, uuid),
	)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClientChecksum", err, 1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = internal.CustomizeError("sendHttpClientResponseChecksum", errors.New("Checksum response error"))
		internal.PrintErrorWithExit("sendHttpClientResponseChecksum", err, 1)
	}
	checksumBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClientReadChecksum", err, 1)
	}
	return string(checksumBytes)
}

func GetDownload(localPort int, fileNameOpt string, uuid string) string {
	reader := bufio.NewReader(os.Stdin)
	client := http.Client{}

	resp, err := client.Post(
		fmt.Sprintf("http://localhost:%v/%v", localPort, uuid),
		"binary/octet-stream",
		nil,
	)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClient", err, 1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code:", resp.StatusCode)
		err = internal.CustomizeError("sendHttpClientResponse", errors.New("Download response error"))
		internal.PrintErrorWithExit("sendHttpClientResponse", err, 1)
	}

	// Set filename from header or from output
	fileName := resp.Header.Get("X-Attachment-Name")
	if fileName == "" {
		internal.PrintErrorWithExit("sendHttpClientInvalidAttachmentName", err, 1)
	}
	if fileNameOpt != "" {
		fileName = fileNameOpt
	}
	fmt.Println("File found:", fileName, fmt.Sprintf("[%v]", resp.Header.Get("X-Mime-Type")))

	// Ask client if they wanna download it or not
	fmt.Print("Download? [(Y)es/(N)o] [Default: Y] -> ")
	prompt, err := reader.ReadString('\n')
	if err != nil {
		internal.PrintErrorWithExit("sendPromptError", err, 1)
	}
	prompt = strings.Replace(prompt, "\n", "", -1)
	if strings.ToLower(prompt) == "n" {
		internal.PrintErrorWithExit("sendPromptCancel", errors.New("Canceled by action"), 0)
	}

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

	return filePath
}
