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

type ConfigSourceAuth struct {
	Default           int          `json:"default"`
	ListConfiguration []ConfigFile `json:"listConfig"`
}

type ConfigFile struct {
	Name  string `json:"name"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Proxy string `json:"proxy"`
	Key   string `json:"key"`
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

func (c ConfigFile) WriteRawConfig() (conf string, err error) {
	var config ConfigSourceAuth

	if CheckConfigFileExist() {
		err = config.ParseConfig(&config)
		if err != nil {
			err = os.Remove(GetConfigPath())
			if err != nil {
				return "", err
			}
		}
	}

	// Check if there's multiple configuration
	if config.GetTunnelBasedInstanceName(c.Name) != -1 {
		return "", errors.New("Multiple instance name")
	}

	// Append config file
	config.ListConfiguration = append(config.ListConfiguration, c)

	strJsonByte, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	conf = base64.StdEncoding.EncodeToString(strJsonByte)
	return conf, nil
}

func (c ConfigFile) WriteConfig() (err error) {
	fmt.Println("Writing config file...")
	strConfig, err := c.WriteRawConfig()
	if err != nil {
		return err
	}

	err = os.WriteFile(ConfigFileLocation, []byte(strConfig), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Config file is saved on", ConfigFileLocation)
	return nil
}

func (c ConfigSourceAuth) WriteConfig() (err error) {
	strJsonByte, err := json.Marshal(c)
	if err != nil {
		return err
	}
	strConfig := base64.StdEncoding.EncodeToString(strJsonByte)
	err = os.WriteFile(ConfigFileLocation, []byte(strConfig), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Config file is saved on", ConfigFileLocation)
	return nil
}

func (c ConfigSourceAuth) ParseConfig(conf *ConfigSourceAuth) error {
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

func (c *ConfigSourceAuth) SetDefault(instanceName string) (err error) {
	instanceID := c.GetTunnelBasedInstanceName(instanceName)
	if instanceID == -1 {
		return errors.New("Instance not found.")
	}

	c.Default = instanceID

	strJsonByte, err := json.Marshal(c)
	if err != nil {
		return err
	}
	strConfig := base64.StdEncoding.EncodeToString(strJsonByte)
	err = os.WriteFile(ConfigFileLocation, []byte(strConfig), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Config file is saved on", ConfigFileLocation)
	return nil
}

func (c ConfigSourceAuth) GetTunnelBasedInstanceName(instanceName string) (index int) {
	index = -1

	for i, v := range c.ListConfiguration {
		if v.Name == instanceName {
			index = i
			break
		}
	}

	return index
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
