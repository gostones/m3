#!/usr/bin/env sh

# https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl/41366949#41366949

openssl req \
    -newkey rsa:2048 \
    -x509 \
    -sha256 \
    -nodes \
    -keyout cert.key \
    -out cert.crt \
    -subj /C=US/ST=CA/L=Bayarea/O=DHNT/OU=M3/CN=m3 \
    -reqexts SAN \
    -extensions SAN \
    -config <(cat openssl.cnf <(printf '[SAN]\nsubjectAltName=DNS:m3,DNS:local.m3,DNS:*.local.m3,DNS:home,DNS:*.home,DNS:*.home.m3')) \
    -days 3650

echo "Done!"
##