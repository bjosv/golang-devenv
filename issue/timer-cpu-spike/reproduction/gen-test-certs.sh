#!/bin/bash -e

# Generate CA
openssl genrsa -out ca.key 4096
openssl req \
        -x509 -new -nodes -sha256 \
        -key ca.key \
        -days 365 \
        -subj "/O=Timer Test/CN=CertAuth Test" \
        -out ca.crt

# Generate cert
openssl genrsa -out test.key 2048
openssl req \
        -new -sha256 \
        -subj "/O=Timer Test/CN=Test" \
        -key test.key | \
            openssl x509 \
                -req -sha256 \
                -CA ca.crt \
                -CAkey ca.key \
                -CAserial ca.txt \
                -CAcreateserial \
                -days 365 \
                -out test.crt
