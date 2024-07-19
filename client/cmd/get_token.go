package main

import (
	"encoding/json"
	"net/http"
)

type CreateRoomJSONReturn struct {
	Token string `json:"token"`

	// This respon fired when the API is failed
	ErrorTitle string `json:"title"`
	Errors     struct {
		Port []string `json:"port"`
	} `json:"errors"`
}

func getToken(path string) (CreateRoomJSONReturn, error) {
	var client = &http.Client{}
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return CreateRoomJSONReturn{} ,err
	}
	res, err := client.Do(req)
	if err != nil {
		return CreateRoomJSONReturn{} ,err
	}
	defer res.Body.Close()

	var roomData CreateRoomJSONReturn
	err = json.NewDecoder(res.Body).Decode(&roomData)
	if err != nil {
		return CreateRoomJSONReturn{} ,err
	}

	return roomData, nil
}
