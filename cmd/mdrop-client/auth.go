package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mplus-oss/mdrop/internal"
)

func AuthCommand() {
	errChan := make(chan error, 0)

	// Set prompt
	config, err := authPrompt()
	if err != nil {
		internal.PrintErrorWithExit("authPromptError", err, 1)
	}

	// Connect to tunnel
	fmt.Println("Connecting to tunnel...")
	sshConfig := internal.SSHParameter{
		ConfigFile: config,
		Command:    []string{"ping"},
	}
	go func() {
		sshOutput, err := StartShellTunnel(sshConfig, false)
		if err != nil {
			errChan <- err
		}
		if strings.Contains(sshOutput, "Pong!") {
			errChan <- nil
		} else {
			errChan <- errors.New("Unexpected output from Tunnel")
		}
	}()
	err = <-errChan
	if err != nil {
		internal.PrintErrorWithExit("authTunnelError", err, 1)
	}

	// Write config file
	err = config.WriteConfig()
	if err != nil {
		internal.PrintErrorWithExit("authWriteConfig", err, 1)
	}
}

func authPrompt() (config internal.ConfigFile, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Instance Name: ")
	instanceName, err := reader.ReadString('\n')
	if err != nil {
		return internal.ConfigFile{}, err
	}
	instanceName = strings.Replace(instanceName, "\n", "", -1)

	fmt.Print("Hostname: ")
	hostname, err := reader.ReadString('\n')
	if err != nil {
		return internal.ConfigFile{}, err
	}
	hostname = strings.Replace(hostname, "\n", "", -1)

	fmt.Print("Port [22]: ")
	port, err := reader.ReadString('\n')
	if err != nil {
		return internal.ConfigFile{}, err
	}
	port = strings.Replace(port, "\n", "", -1)
	if port == "" {
		port = "22"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return internal.ConfigFile{}, err
	}

	fmt.Print("Proxy [Set blank if none]: ")
	proxy, err := reader.ReadString('\n')
	if err != nil {
		return internal.ConfigFile{}, err
	}
	proxy = strings.Replace(proxy, "\n", "", -1)

	config = internal.ConfigFile{
		Name:  instanceName,
		Host:  hostname,
		Port:  portInt,
		Proxy: proxy,
	}
	return config, nil
}
