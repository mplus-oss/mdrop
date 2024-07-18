#!/bin/sh

# Generate RS2048A Private Key
openssl genrsa -out prikey.pem 2048

# Generate Public Key
openssl rsa -in prikey.pem -pubout -out pubkey.pem

# Generate Certificate for 1 year
openssl req -new -x509 -key prikey.pem -out cert.pem -days 365 \
  -subj "/CN=mplus.software/OU=Mplus Software/O=Mplus DevOps Team/L=South Jakarta/ST=Greater Jakarta/C=ID"

# Move the cert.pem on development folder
mv *pem ./MDrop.Broker
