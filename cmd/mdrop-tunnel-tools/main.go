package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: mdrop-tunnel <subcommand> <args>")
		os.Exit(1)
	}

	cmd := args[0]
	switch cmd {
	case "ping":
		PingCommand()
	case "get-port":
		GetPortCommand()
	case "connect":
		ConnectCommand()
	default:
		fmt.Println("Subcommand: ping, get-port")
	}
}

func ConnectCommand() {
	for true {
		fmt.Println("Pong.....")
		time.Sleep(10 * time.Second)
	}
}

func GetPortCommand() {
	portFound := 0
	for i := 5000; i <= 59999; i++ {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(i)), timeout)
		if err != nil {
			fmt.Println("Error ["+strconv.Itoa(i)+"]: "+err.Error())
		}
		if conn == nil {
			portFound = i
			break
		}
		fmt.Println("Closed ["+strconv.Itoa(i)+"]: Continue.")
	}

	if portFound != 0 {
		fmt.Println("Get-Port-Found "+strconv.Itoa(portFound))
	} else {
		fmt.Println("Closed ["+strconv.Itoa(portFound)+"]: Not found.")
		os.Exit(2)
	}
}

func PingCommand() {
	fmt.Println("Pong!")
}
