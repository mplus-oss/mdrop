package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type SSHParameter struct {
	IsRemote	bool
	LocalPort	int
	RemotePort	int
	ConfigFile	ConfigFile
	Command		[]string
}

func (s SSHParameter) GenerateConnectSSHArgs() string {
	args := []string{
		"-p",
		strconv.Itoa(s.ConfigFile.Port),
		"-o",
		"UserKnownHostsFile=/dev/null",
		"-o",
		"StrictHostKeyChecking=no",
		"tunnel@" + s.ConfigFile.Host,
	}
	args = append(args, s.Command...)
	s.GenerateProxyArgs(&args)

	return "ssh "+strings.Join(args, " ")
}

func (s SSHParameter) GenerateRemoteSSHArgs() string {
	args := []string{
		"-p",
		strconv.Itoa(s.ConfigFile.Port),
		"-o",
		"UserKnownHostsFile=/dev/null",
		"-o",
		"StrictHostKeyChecking=no",
		"tunnel@" + s.ConfigFile.Host,
		"connect",
	}

	if s.IsRemote {
		args = append(
			[]string{"-R", fmt.Sprintf("%v:127.0.0.1:%v", s.RemotePort, s.LocalPort)},
			args...,
		)
	} else {
		args = append(
			[]string{"-L", fmt.Sprintf("%v:127.0.0.1:%v", s.LocalPort, s.RemotePort)},
			args...,
		)
	}
	s.GenerateProxyArgs(&args)

	return "ssh "+strings.Join(args, " ")
}

func (s SSHParameter) GenerateProxyArgs(args *[]string) {
	switch s.ConfigFile.Proxy {
	case "cloudflared":
		*args = append(
			[]string{
				"-o",
				fmt.Sprintf("ProxyCommand=\"cloudflared access ssh --hostname %v\"", s.ConfigFile.Host),
			},
			*args...,
		)
	}
}
