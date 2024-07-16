#!/usr/bin/env bash

args="$(echo "$SSH_ORIGINAL_COMMAND" | awk -F' ' '{print NF}')"

if [ $args != "1" ]; then
	echo "WRONG_USAGE_EXCEPTION"
	exit 1
fi

# MOTD
port="$(echo $SSH_ORIGINAL_COMMAND | awk '{ print $1 }')"
echo "Connected! $port"

tail -f /dev/null
