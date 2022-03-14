#!/bin/bash -e

generate_cert() {
    local dir="$1"
    local name="$2"
    local ca="$3"
    local cn="$4"

    local keyfile=${dir}/${name}.key
    local certfile=${dir}/${name}.crt

    echo "*********************************************************"
    echo "Generate keypair for $name"
    echo "*********************************************************"
    echo

    # Set domain for subjectAltName (DNS or IP)
    local domain="DNS"
    if [[ $cn =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        domain="IP"
    fi

    openssl genrsa -out $keyfile 2048
    openssl req \
        -new -sha256 \
        -subj "/O=Redis Test/CN=$cn" \
        -key $keyfile | \
            openssl x509 \
                -req -sha256 \
                -CA ${dir}/${ca}.crt \
                -CAkey ${dir}/${ca}.key \
                -CAserial ${dir}/${ca}.txt \
                -CAcreateserial \
                -days 365 \
                -extfile <(printf "[ EXT ]
                                  keyUsage = digitalSignature, keyEncipherment
                                  subjectAltName = $domain:$cn") \
                -extensions EXT \
                -out $certfile
}

generate_ca() {
    local dir="$1"
    local name="$2"
    echo "*********************************************************"
    echo "Generate CA ${name}"
    echo "*********************************************************"
    echo

    openssl genrsa -out ${dir}/${name}.key 4096
    openssl req \
        -x509 -new -nodes -sha256 \
        -key ${dir}/${name}.key \
        -days 365 \
        -subj "/O=Redis Test (${name})/CN=CertAuth ${name}" \
        -out ${dir}/${name}.crt
}


# Generates TLS certificates to:
dir=/tmp/tls-data

ca='ca'
cn='localhost'

mkdir -p $dir
generate_ca $dir $ca

generate_cert $dir "redis"    $ca $cn
generate_cert $dir "exporter" $ca $cn

# Let the pods read the key files
chmod 644 $dir/*.key
