#!/bin/sh
set -e

# Set PWD to /app
cd /app

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

# Launch the app
env X509_CERTIFICATE_PASSWORD="$RANDOM_STRING" /app/MDrop.Broker
