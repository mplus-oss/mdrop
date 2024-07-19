package internal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ConfigFileLocation string
)

type ConfigFile struct {
	Token	string
	URL		string
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

	strConfig := base64.StdEncoding.EncodeToString([]byte(c.URL))
	if c.Token != "" {
		strConfig += " " + base64.RawStdEncoding.EncodeToString([]byte(c.Token))
	}
	err := os.WriteFile(ConfigFileLocation, []byte(strConfig), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Config file is saved on", ConfigFileLocation)
	return nil
}

func (c ConfigFile) ParseConfig(conf* ConfigFile) (error) {
	var urlByte, tokenByte []byte

	if !CheckConfigFileExist() {
		return errors.New("No config file on local. Please log in first.")
	}
	file, err := os.ReadFile(ConfigFileLocation)
	if err != nil {
		return err
	}

	fileSplit := strings.Split(string(file), " ")
	if len(fileSplit) == 0 && len(fileSplit) > 2 {
		return errors.New("Invalid config file. Please log in again.")
	}

	urlByte, err = base64.StdEncoding.DecodeString(fileSplit[0])
	if len(fileSplit) > 1 {
		tokenByte, err = base64.RawStdEncoding.DecodeString(fileSplit[1])
	} else {
		tokenByte = []byte("")
	}
	if err != nil {
		return err
	}

	conf.URL = string(urlByte)
	conf.Token = string(tokenByte)

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
