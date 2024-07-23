package internal

import (
	"fmt"
)

func GenerateSSHArgs(isRemote bool, c ConfigFile, localPort int, remotePort int) string {
	args := fmt.Sprintf(
		"-p %v -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no tunnel@%v",
		c.Tunnel.Port,
		c.Tunnel.Host,
	)
	
	flag := [4]int{}
	remoteFlag := ""
	flag[2] = remotePort
	if isRemote {
		remoteFlag = "-R"
		flag[0] = remotePort
		flag[1] = localPort
	} else {
		remoteFlag = "-L"
		flag[1] = remotePort
		flag[0] = localPort
	}

	args = fmt.Sprintf(
		"%v %v:127.0.0.1:%v %v %v",
		remoteFlag,
		flag[0],
		flag[1],
		args,
		flag[2],
	)
	return args
}
