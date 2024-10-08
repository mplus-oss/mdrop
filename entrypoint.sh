#!/bin/bash

set -e

# Create keyfile private for mdrop-tunnel
touch /etc/mdrop-tunnel-key
if [[ "$PRIVATE_MODE_TOKEN" != "" ]]; then
	echo "$PRIVATE_MODE_TOKEN" >> /etc/mdrop-tunnel-key
fi

# Run MOTD on logs
echo "SSHD successfully launched. For reference, you can use this command to remote port to the container:"
echo "    ssh -R <remote>:127.0.0.1:<local> -T -p <port> -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no tunnel@localhost <command>"
echo "Or listening port from the container:"
echo "    ssh -L <local>:127.0.0.1:<remote> -T -p <port> -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no tunnel@localhost <command>"

# Run and detach SSHD
/usr/sbin/sshd -D
