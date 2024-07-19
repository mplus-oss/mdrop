package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func AuthCommand(args []string) {
	flag := flag.NewFlagSet("mdrop auth", flag.ExitOnError)
	var (
		url = flag.String("url", "https://example.com", "URL of broker")
		token = flag.String("token", "token-is-here", "Private key token if the broker is on private mode")
	)
	flag.Parse(args)

	var client = &http.Client{}

	req, err := http.NewRequest("POST", *url + "/verify?token=" + *token, nil)
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

	var verifyData VerifyJSONReturn
	err = json.NewDecoder(res.Body).Decode(&verifyData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		var msg string
		if verifyData.ErrorTitle == "" {
			msg = verifyData.Message
		} else {
			msg = verifyData.ErrorTitle
		}
		fmt.Println("Error:", msg)
		os.Exit(1)
	}

	base64Url := base64.StdEncoding.EncodeToString([]byte(*url))
	base64Token := ""
	if verifyData.IsPublic {
		fmt.Println("This broker is running on public mode.")
	} else {
		base64Token = base64.StdEncoding.EncodeToString([]byte(*token))
	}

	fmt.Println("Authenticated!")
	fmt.Println(base64Url)
	fmt.Println(base64Token)
	os.Exit(0)
}

type VerifyJSONReturn struct {
	Message		 string `json:"message"`
	IsPublic	 bool	`json:"isPublic"`

	// This respon fired when the API is failed
	ErrorTitle	 string `json:"title"`
}
