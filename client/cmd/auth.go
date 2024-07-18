package main

import (
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

	fmt.Println(res.StatusCode)
	fmt.Println(verifyData.Message)
}

type VerifyJSONReturn struct {
	Message string `json:"message"`
}
