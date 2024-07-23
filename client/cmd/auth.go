package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	internal "github.com/mplus-oss/mdrop/client/internal"
)

func AuthCommand(args []string) {
	flag := flag.NewFlagSet("mdrop auth", flag.ExitOnError)
	var (
		url   = flag.String("url", "https://example.com", "URL of broker")
		token = flag.String("token", "", "Private key token if the broker is on private mode")
	)
	flag.Parse(args)

	if *token == "" {
		*token = "token"
	}

	var client = &http.Client{}

	fmt.Println("Authenticating...")
	req, err := http.NewRequest("POST", *url+"/verify?token="+*token, nil)
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

	if verifyData.IsPublic {
		fmt.Println("This broker is running on public mode.")
		*token = ""
	}

	config := internal.ConfigFile{
		URL:    *url,
		Token:  *token,
		Tunnel: internal.ConfigFileTunnel{
			Host: verifyData.TunnelHost,
			Port: verifyData.TunnelPort,
		},
	}
	err = config.WriteConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Authenticated!")
	os.Exit(0)
}

type VerifyJSONReturn struct {
	Message		 string `json:"message"`
	IsPublic	 bool   `json:"isPublic"`
	TunnelHost	 string `json:"tunnelHost"`
	TunnelPort	 int    `json:"tunnelPort"`

	// This respon fired when the API is failed
	ErrorTitle string `json:"title"`
}
