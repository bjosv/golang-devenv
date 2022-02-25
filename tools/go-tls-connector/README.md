# Go TLS test

## Prepare TLS certs
./gen-test-certs.sh
./gen-test-certs.sh redis 600 db

## Build test container
docker build -t bjosv/go-tls-connector:1.16.6 -f Dockerfile.go1.16.6 .
docker build -t bjosv/go-tls-connector:1.16.7 -f Dockerfile.go1.16.7 .

## Run
docker-compose up
