#!/usr/bin/env sh

## generate cert for my peer id

filename="cert-myid"

myid=$(hexid -id $(ipfs id -f="<id>\n")); if [ $? -ne 0 ]; then
    echo "Failed to get my peer id."
    exit 1
fi

openssl req \
    -newkey rsa:2048 \
    -x509 \
    -sha256 \
    -nodes \
    -keyout ${filename}.key \
    -out ${filename}.crt \
    -subj /C=US/ST=CA/L=Bayarea/O=DHNT/OU=M3/CN=${myid}.m3 \
    -reqexts SAN \
    -extensions SAN \
    -config <(cat openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:${myid}.m3,DNS:*.${myid}.m3")) \
    -days 3650

echo "Done!"
##