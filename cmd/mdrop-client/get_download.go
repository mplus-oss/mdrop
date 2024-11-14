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

func GetPrompt(localPort int, uuid string, fileNameOpt string, isSingleFile bool, isSkipPrompt bool) (filePath string) {
	reader := bufio.NewReader(os.Stdin)
	client := http.Client{}

	resp, err := client.Get(
		fmt.Sprintf("http://localhost:%v/verify-%v", localPort, uuid),
	)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClientPrompt", err, 1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = internal.CustomizeError("sendHttpClientResponsePrompt", errors.New("Prompt response error"))
		internal.PrintErrorWithExit("sendHttpClientResponsePrompt", err, 1)
	}
	fileNameBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		internal.PrintErrorWithExit("sendHttpClientReadPrompt", err, 1)
	}
	fileName := string(fileNameBytes)

	// Ask client if they wanna download it or not
	fmt.Println("\nFile found:", fileName, fmt.Sprintf("[%v]", resp.Header.Get("X-Mime-Type")))
    if !isSkipPrompt {
        fmt.Print("Download? [(Y)es/(N)o] [Default: Y] -> ")
        prompt, err := reader.ReadString('\n')
        if err != nil {
            internal.PrintErrorWithExit("sendPromptError", err, 1)
        }
        prompt = strings.Replace(prompt, "\n", "", -1)
        if strings.ToLower(prompt) == "n" {
            internal.PrintErrorWithExit("sendPromptCancel", errors.New("Canceled by action"), 0)
        }
    }

	// Check if there's duplicate file
	filePath, err = os.Getwd()
	if err != nil {
		internal.PrintErrorWithExit("sendFileWorkDir", err, 1)
	}
	// Check if it's single file
	if isSingleFile {
		if fileNameOpt != "" {
			fmt.Println("Changing filename to", fileNameOpt)
			fileName = fileNameOpt
		}
	}
	if fileStatus, _ := os.Stat(filePath + "/" + fileName); fileStatus != nil {
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
			fmt.Print("Change filename [" + fileName + "]: ")
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

	return fileName
}

func GetDownload(localPort int, uuid string, fileName string, fileNameOpt string, isSingleFile bool) string {
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
	filePath, err := os.Getwd()
	if err != nil {
		internal.PrintErrorWithExit("sendFileWorkDir", err, 1)
	}
	if !isSingleFile {
        if fileNameOpt != "" {
            filePath = fileNameOpt
            fmt.Println("Changing working directory to", filePath)
        }
	}
	filePath += "/" + fileName

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		internal.PrintErrorWithExit("sendFileCreation", err, 1)
	}

	// Downloading file
	fmt.Println("\nDownloading...")
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
