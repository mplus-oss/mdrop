package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/mplus-oss/mdrop/internal"
	"golang.org/x/term"
)

func AuthCommand(args []string) {
	errChan := make(chan error, 0)

	flag := flag.NewFlagSet("mdrop auth", flag.ExitOnError)
	var (
		defaultInstance = flag.String("set-default", "", "Set default tunnel instance.")
		list            = flag.Bool("list", false, "Get list of tunnel instance")
		help            = flag.Bool("help", false, "Print this message")
		deleteInstance  = flag.String("delete", "", "Delete tunnel instance")
		raw             = flag.Bool("raw", false, "Print raw output of .mdrop file")
		replace         = flag.Bool("replace", false, "Replace .mdrop file with new config")
	)
	flag.Parse(args)

	if *help {
		fmt.Println("Command: mdrop auth [options]")
		flag.Usage()
		os.Exit(1)
	}

	if *list {
		getTunnelInstance()
		os.Exit(0)
	}

	if *defaultInstance != "" {
		setDefaultTunnelInstance(*defaultInstance)
		os.Exit(0)
	}

	if *deleteInstance != "" {
		deleteTunnelInstance(*deleteInstance)
		os.Exit(0)
	}

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
	if *raw {
		confString, err := config.WriteRawConfig(*replace)
		if err != nil {
			internal.PrintErrorWithExit("authWriteRawConfig", err, 1)
		}
		fmt.Println("\nCopy this secret below and paste it to ~/.mdrop")
		fmt.Println(confString)
		os.Exit(0)
	}
	err = config.WriteConfig(*replace)
	if err != nil {
		internal.PrintErrorWithExit("authWriteConfig", err, 1)
	}
}

func setDefaultTunnelInstance(instanceName string) {
	var authConfig internal.ConfigSourceAuth
	err := authConfig.ParseConfig(&authConfig)
	if err != nil {
		internal.PrintErrorWithExit("authParseConfig", err, 1)
	}

	err = authConfig.SetDefault(instanceName)
	if err != nil {
		internal.PrintErrorWithExit("authSetDefaultTunnelInstance", err, 1)
	}
}

func getTunnelInstance() {
	var authConfig internal.ConfigSourceAuth
	err := authConfig.ParseConfig(&authConfig)
	if err != nil {
		internal.PrintErrorWithExit("authParseConfig", err, 1)
	}

	fmt.Println("Default instance:", authConfig.ListConfiguration[authConfig.Default].Name)
	for i, instance := range authConfig.ListConfiguration {
		fmt.Println(
			fmt.Sprintf("[%v] %v", i, instance.Name),
		)
	}
}

func deleteTunnelInstance(instanceName string) {
	var authConfig internal.ConfigSourceAuth
	err := authConfig.ParseConfig(&authConfig)
	if err != nil {
		internal.PrintErrorWithExit("authParseConfig", err, 1)
	}

	instanceIndex := -1
	for i, instance := range authConfig.ListConfiguration {
		if instance.Name == instanceName {
			instanceIndex = i
			break
		}
	}
	if instanceIndex < 0 {
		internal.PrintErrorWithExit("authInstanceNotFound", errors.New("Instance not found"), 1)
	}
	authConfig.ListConfiguration = append(
		authConfig.ListConfiguration[:instanceIndex],
		authConfig.ListConfiguration[instanceIndex+1:]...,
	)
	err = authConfig.WriteConfig()
	if err != nil {
		internal.PrintErrorWithExit("authParseConfig", err, 1)
	}
}

func authPrompt() (config internal.ConfigFile, err error) {
	reader := bufio.NewReader(os.Stdin)
	stdin := int(syscall.Stdin)
	oldState, err := term.GetState(stdin)
	if err != nil {
		return internal.ConfigFile{}, err
	}
	defer term.Restore(stdin, oldState)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		for _ = range sigch {
			term.Restore(stdin, oldState)
			os.Exit(1)
		}
	}()

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

	fmt.Print("Private Key String [Set blank if none]: ")
	privateKeyByte, err := term.ReadPassword(stdin)
	if err != nil {
		return internal.ConfigFile{}, err
	}
	privateKeyString := strings.Replace(string(privateKeyByte), "\n", "", -1)

	config = internal.ConfigFile{
		Name:  instanceName,
		Host:  hostname,
		Port:  portInt,
		Proxy: proxy,
		Key:   privateKeyString,
	}
	return config, nil
}
