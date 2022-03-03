# Go development environment

## Install multiple versions
go install golang.org/dl/go1.16.6@latest
go1.16.6 download
go1.16.6 version

go install golang.org/dl/go1.16.7@latest
go1.16.7 download
go1.16.7 version

go install golang.org/dl/go1.18beta1@latest
go1.18beta1 download

## Go base containers

### Docker Official Image packaging for golang
https://github.com/docker-library/golang.git

#### Docker file for 1.16.7-alpine aka 1.16.7-alpine3.14
https://github.com/docker-library/golang/blob/4c1da70f967b2b38b254e166e787d017cc9ca351/1.16/alpine3.14/Dockerfile

#### Docker file for 1.16.6-alpine aka 1.16.6-alpine3.14
https://github.com/docker-library/golang/blob/54aa949c354b1e14cb636539f401b0e58ca76927/1.16/alpine3.14/Dockerfile

Build jobs: https://github.com/docker-library/golang/runs/3254891813?check_suite_focus=true

https://github.com/docker-library/docs/commit/14ed5488194a3320b9a5c5f1f09df3585e3fcb28
[`1.16.7-alpine3.14`, `1.16-alpine3.14`, `1-alpine3.14`, `alpine3.14`, `1.16.7-alpine`, `1.16-alpine`, `1-alpine`, `alpine`]


## Build own base container
docker build -t golang:1.16.6-alpine-own images/golang1.16.6-alpine
docker build -t golang:1.16.7-alpine-own images/golang1.16.7-alpine

cd images/golang-master
git clone https://github.com/golang/go.git
docker build -t golang:master --no-cache images/golang-master
