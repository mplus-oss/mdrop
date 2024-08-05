#!/bin/bash

set -e

# Set PWD to /root
cd /

# Set constant
RANDOM_STRING=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 64; echo)

# Generate RS2048A Private Key
echo "Generate Private Key"
openssl genrsa -passout "pass:$RANDOM_STRING" -aes256 -out prikey.pem 2048

# Generate Public Key
echo "Generate Public Key"
openssl rsa -in prikey.pem -passin "pass:$RANDOM_STRING" -passout "pass:$RANDOM_STRING" -pubout -out pubkey.pem

# Generate Certificate for 1 year
echo "Generate Certificate"
openssl req -new -x509 -key prikey.pem -passin "pass:$RANDOM_STRING" -passout "pass:$RANDOM_STRING" -out cert.pem -days 365 \
  -subj "/CN=mplus.software/OU=Mplus Software/O=Mplus DevOps Team/L=South Jakarta/ST=Greater Jakarta/C=ID"

# Print all env to file
echo "$RANDOM_STRING" >> /.cert
echo "$PRIVATE_TOKEN" >> /.token

# Run MOTD on logs
echo "SSHD successfully launched. For reference, you can use this command to remote port to the container:"
echo "    ssh -R <remote>:127.0.0.1:<local> -T -p <port> -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no tunnel@localhost <command>"
echo "Or listening port from the container:"
echo "    ssh -L <local>:127.0.0.1:<remote> -T -p <port> -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no tunnel@localhost <command>"

# Run and detach SSHD
/usr/sbin/sshd -D
