package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	internal "github.com/mplus-oss/mdrop/client/internal"
)

func SendCommand(args []string) {
	flag := flag.NewFlagSet("mdrop send", flag.ExitOnError)
	var (
		token = flag.String("token", "", "Token from reciever")
		file  = flag.String("file", "", "File that want to send to reciever")
	)
	flag.Parse(args)

	if *token == "" {
		fmt.Println("No token provided!")
		os.Exit(1)
	}
	if *file == "" {
		fmt.Println("No file provided!")
		os.Exit(1)
	}

	var c internal.ConfigFile
	err := c.ParseConfig(&c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Decrypt token first
	client := &http.Client{}
	path := fmt.Sprintf(
		"%v/room/join?token=%v",
		c.URL,
		*token,
	)
	fmt.Println("Authenticating...")
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	var joinData JoinJSONReturn
	err = json.NewDecoder(res.Body).Decode(&joinData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		var msg string
		if joinData.ErrorTitle == "" {
			msg = joinData.Message
		} else {
			msg = joinData.ErrorTitle
		}
		fmt.Println("Error:", msg)
		os.Exit(1)
	}

	// Try to send data to reciever
	fmt.Println("Authenticated. Connecting to reciever...")
	fmt.Println(joinData.Port)
}

func uploadFile(uri string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", "")
	return req, err
}

type JoinJSONReturn struct {
	Message  string `json:"message"`
	Port	 int	`json:"port"`

	// This respon fired when the API is failed
	ErrorTitle string `json:"title"`
}
