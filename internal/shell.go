package internal

import (
	"fmt"
	"strconv"
	"strings"
)

func GenerateSSHArgs(isRemote bool, c ConfigFile, localPort int, remotePort int) string {
	args := []string{
		"-o",
		"UserKnownHostsFile=/dev/null",
		"-o",
		"StrictHostKeyChecking=no",
		"tunnel@" + c.Tunnel.Host,
		strconv.Itoa(remotePort),
	}
	generateProxyCommand(&args, c)

	if isRemote {
		args = append(
			[]string{"-R", fmt.Sprintf("%v:127.0.0.1:%v", remotePort, localPort)},
			args...,
		)
	} else {
		args = append(
			[]string{"-L", fmt.Sprintf("%v:127.0.0.1:%v", localPort, remotePort)},
			args...,
		)
	}

	return "ssh " + strings.Join(args, " ")
}

func generateProxyCommand(args *[]string, c ConfigFile) {
	switch c.Tunnel.Proxy {
	case "cloudflared":
		*args = append(
			[]string{
				"-o",
				fmt.Sprintf("ProxyCommand=\"cloudflared access ssh --hostname %v\"", c.Tunnel.Host),
			},
			*args...,
		)
	default:
		*args = append([]string{"-p", strconv.Itoa(c.Tunnel.Port)}, *args...)
	}
}
