#1/bin/sh

# Run MOTD on logs
echo "SSHD successfully launched. For reference, you can use this command to remote port to the container:"
echo "    ssh -R <remote>:127.0.0.1:<local> -p <port> -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no tunnel@localhost <remote>"
echo "Or listening port from the container:"
echo "    ssh -L <local>:127.0.0.1:<remote> -p <port> -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no tunnel@localhost <remote>"

# Run and detach SSHD
/usr/sbin/sshd -D
