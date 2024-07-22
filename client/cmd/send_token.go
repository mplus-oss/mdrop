package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	internal "github.com/mplus-oss/mdrop/client/internal"
)

type JoinRoomJSONReturn struct {
	Message  string `json:"message"`
	Port	 int	`json:"port"`

	// This respon fired when the API is failed
	ErrorTitle string `json:"title"`
	Errors     struct {
		Port []string `json:"port"`
	} `json:"errors"`
}

func sendToken(c internal.ConfigFile, token string) (JoinRoomJSONReturn, error) {
	// Decrypt token first
	client := &http.Client{}
	path := fmt.Sprintf(
		"%v/room/join?token=%v",
		c.URL,
		token,
	)
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

	var joinData JoinRoomJSONReturn
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

	return joinData, err
}
