package main

import (
	"encoding/base64"
	"encoding/json"
)

type TokenTransferJSON struct {
	Host       string   `json:"host"`
	RemotePort int      `json:"remotePort"`
	Files      []string `json:"files"`
}

func (t TokenTransferJSON) GenerateToken() (token string, err error) {
	jsonByte, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(jsonByte), nil
}

func (t TokenTransferJSON) ParseToken(token string, tokenConfig *TokenTransferJSON) (err error) {
	jsonByte, err := base64.RawStdEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonByte, &tokenConfig)
	if err != nil {
		return err
	}
	return nil
}
