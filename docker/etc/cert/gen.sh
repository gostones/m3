#!/usr/bin/env bash

cat << EOF > cert.cnf
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
countryName = US
countryName_default = US
stateOrProvinceName = CA
stateOrProvinceName_default = CA
localityName = Bay Area
localityName_default = Bay Area
organizationalUnitName = M3
organizationalUnitName_default = M3
commonName = *.home
commonName_max = 64

[v3_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = home
DNS.2 = *.home
EOF

openssl genrsa -out cert.key 2048
openssl req -new -out cert.csr -key cert.key -config cert.cnf
openssl x509 -req -days 3650 -in cert.csr -signkey cert.key -out cert.pem -extensions v3_req -extfile cert.cnf

#---
# #
# SSL_DIR="./etc/ssl/xip.io"

# # Set the wildcarded domain
# # we want to use
# DOMAIN="*.xip.io"

# # A blank passphrase
# PASSPHRASE=""

# # Set our CSR variables
# SUBJ="
# C=US
# ST=Connecticut
# O=
# localityName=New Haven
# commonName=$DOMAIN
# organizationalUnitName=
# emailAddress=
# "

# # Create our SSL directory
# # in case it doesn't exist
# sudo mkdir -p "$SSL_DIR"

# # Generate our Private Key, CSR and Certificate
# sudo openssl genrsa -out "$SSL_DIR/xip.io.key" 2048
# sudo openssl req -new -subj "$(echo -n "$SUBJ" | tr "\n" "/")" -key "$SSL_DIR/xip.io.key" -out "$SSL_DIR/xip.io.csr" -passin pass:$PASSPHRASE
# sudo openssl x509 -req -days 365 -in "$SSL_DIR/xip.io.csr" -signkey "$SSL_DIR/xip.io.key" -out "$SSL_DIR/xip.io.crt"





# #https://ksearch.wordpress.com/2017/08/22/generate-and-import-a-self-signed-ssl-certificate-on-mac-osx-sierra/

# openssl genrsa -des3 -passout pass:x -out server.pass.key 2048
# openssl rsa -passin pass:x -in server.pass.key -out server.key
# rm server.pass.key

# # #
# openssl req -new -key server.key -out server.csr

# echo ""
# echo "*** specify-the-same-common-name-that-you-used-while-generating-csr-in-the-last-step ***"
# echo "Common Name (e.g. server FQDN or YOUR name):"
# echo ""
# read FQDN

# cat << EOF > v3.ext
# authorityKeyIdentifier=keyid,issuer
# basicConstraints=CA:FALSE
# keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
# subjectAltName = @alt_names

# [alt_names]
# DNS.1 = $FQDN
# EOF

# openssl x509 -req -sha256 -extfile v3.ext -days 365 -in server.csr -signkey server.key -out server.crt
# openssl x509 -in server.crt -out server.pem -outform PEM

# #
#---
# openssl genrsa -des3 -out hostname.key 2048
# openssl rsa -in hostname.key -out hostname-key.pem
# openssl req -new -key hostname-key.pem -out hostname-request.csr
# openssl x509 -req -extensions v3_req -days 365 -in hostname-request.csr -signkey hostname-key.pem -out hostname-cert.pem -extfile <path to openssl.conf>