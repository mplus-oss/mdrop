package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	ConfigFileLocation string
)

type ConfigFileTunnel struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Proxy string `json:"proxy"`
}

type ConfigFile struct {
	Token  string           `json:"token"`
	URL    string           `json:"url"`
	Tunnel ConfigFileTunnel `json:"tunnel"`
}

func init() {
	var err error

	ConfigFileLocation, err = os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ConfigFileLocation += "/.mdrop"
}

func (c ConfigFile) WriteConfig() error {
	fmt.Println("Writing config file...")
	if CheckConfigFileExist() {
		err := os.Remove(GetConfigPath())
		if err != nil {
			return err
		}
	}

	strJsonByte, err := json.Marshal(c)
	if err != nil {
		return err
	}
	strConfig := base64.StdEncoding.EncodeToString(strJsonByte)
	if c.Token != "" {
		strConfig += " " + base64.RawStdEncoding.EncodeToString([]byte(c.Token))
	}
	err = os.WriteFile(ConfigFileLocation, []byte(strConfig), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Config file is saved on", ConfigFileLocation)
	return nil
}

func (c ConfigFile) ParseConfig(conf *ConfigFile) error {
	if !CheckConfigFileExist() {
		return errors.New("No config file on local. Please log in first.")
	}
	file, err := os.ReadFile(ConfigFileLocation)
	if err != nil {
		return err
	}

	confByte, err := base64.StdEncoding.DecodeString(string(file))
	if err != nil {
		return err
	}
	err = json.Unmarshal(confByte, &conf)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigPath() string {
	return ConfigFileLocation
}

func CheckConfigFileExist() bool {
	if _, err := os.Stat(ConfigFileLocation); err != nil {
		return false
	}
	return true
}
