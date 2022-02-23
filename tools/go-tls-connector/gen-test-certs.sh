#!/bin/bash -e

# Generates TLS certificates to:
dir=/tmp/tls-data

generate_cert() {
    local name="$1"
    local ca="$2"
    local cn="$3"
    local faketime="$4"
    local type="$5"

    local keyfile=${dir}/${name}.key
    local certfile=${dir}/${name}.crt

    echo "*********************************************************"
    echo "Generate keypair for $name using faketime=$faketime"
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
        faketime -f ${faketime} \
            openssl x509 \
                -req -sha256 \
                -CA ${dir}/${ca}.crt \
                -CAkey ${dir}/${ca}.key \
                -CAserial ${dir}/${ca}.txt \
                -CAcreateserial \
                -days 1 \
                -extfile <(printf "[ EXT ]
                                  keyUsage = digitalSignature, keyEncipherment
                                  nsCertType = $type
                                  subjectAltName = $domain:$cn") \
                -extensions EXT \
                -out $certfile
}

generate_ca() {
    local name="$1"
    local faketime="$2"
    echo "*********************************************************"
    echo "Generate CA ${name} using faketime=$faketime"
    echo "*********************************************************"
    echo

    openssl genrsa -out ${dir}/${name}.key 4096
    faketime -f ${faketime} \
        openssl req \
            -x509 -new -nodes -sha256 \
            -key ${dir}/${name}.key \
            -days 1 \
            -subj "/O=Redis Test (${name})/CN=CertAuth ${name}" \
            -out ${dir}/${name}.crt
}


mkdir -p ${dir}

# Arguments:
#   Generated a specific keypair or empty string for all
name=$1
#   Minutes until certs expire, or skip for 1 day until expire
minutes=$2
#   Override default common name i.e the server host name
cn=${3:-'localhost'}

faketime="+0m"

# Create cert's with X minute until expiring (60min * 24h - X min),
# example: 2 give faketime="-1438m"
[[ ! -z $minutes ]] && faketime="-$((1440-$minutes))m"

# Generate CA if no specific keypair name given
ca="ca"
#ca_s="ca_server"
#ca_c="ca_client"
[[ -z $name || "$name" == "ca" ]] && generate_ca $ca $faketime
#[[ -z $name || "$name" == "ca" || "$name" == "ca_server" ]] && generate_ca $ca_s $faketime
#[[ -z $name || "$name" == "ca" || "$name" == "ca_client" ]] && generate_ca $ca_c $faketime
#[[ -z $name || "$name" == ca* ]] && cat ${dir}/${ca_s}.crt ${dir}/${ca_c}.crt > ${dir}/${ca}.crt

# for i in ${ca_s} ${ca_c} ; do
#      openssl x509 -in ${dir}/$i.crt -text >> ${dir}/${ca}.crt
# done

# Generate if no argument, or specific given
# generate_cert <name> <ca-name> <common-name> <faketime> <cert-type>

[[ -z $name || "$name" == "redis" ]]      && generate_cert redis      $ca $cn $faketime "server"
[[ -z $name || "$name" == "exporter-s" ]] && generate_cert exporter-s $ca $cn $faketime "server"
[[ -z $name || "$name" == "exporter-c" ]] && generate_cert exporter-c $ca $cn $faketime "client"
[[ -z $name || "$name" == "curl" ]]       && generate_cert curl       $ca $cn $faketime "client"

# Let the pods read the key files
chmod 644 ${dir}/*.key
