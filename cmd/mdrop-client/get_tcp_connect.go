package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/mplus-oss/mdrop/internal"
)

func GetTcpReadyConnect(localPort int) {
	isConnected := false
	for i := 1; i <= 60; i++ {
		timeout := time.Second
		_, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(localPort)), timeout)
		if err == nil {
			isConnected = true
			break
		}
		time.Sleep(timeout)
	}

	if isConnected {
		fmt.Println("Tunnel is ready!")
	} else {
		internal.PrintErrorWithExit("getTcpReadyConnectError", errors.New("Tunnel is timeout."), 1)
	}
}
