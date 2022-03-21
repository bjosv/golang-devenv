# Go development environment

## Install multiple versions

```
go install golang.org/dl/go1.16.6@latest
go1.16.6 download
go1.16.6 version

# Which uses:
# GOROOT="<HOME>/sdk/go1.16.6"

go install golang.org/dl/go1.16.7@latest
go1.16.7 download

go install golang.org/dl/go1.18beta1@latest
go1.18beta1 download
```

## Go base containers

* Official docker images for golang: https://github.com/docker-library/golang.git

* Info regarding tags: https://github.com/docker-library/docs/tree/master/golang

* Build jobs: https://github.com/docker-library/golang/runs/3254891813?check_suite_focus=true

* [Dockerfile for 1.16.7-alpine aka 1.16.7-alpine3.14](https://github.com/docker-library/golang/blob/4c1da70f967b2b38b254e166e787d017cc9ca351/1.16/alpine3.14/Dockerfile)

* [Dockerfile for 1.16.6-alpine aka 1.16.6-alpine3.14](https://github.com/docker-library/golang/blob/54aa949c354b1e14cb636539f401b0e58ca76927/1.16/alpine3.14/Dockerfile)


## Build own base container

```
docker build -t golang:1.16.6-alpine-own docker/golang1.16.6-alpine
docker build -t golang:1.16.7-alpine-own docker/golang1.16.7-alpine
```

## Build own base container from a go git repo

```
cd docker/golang-git
git clone https://github.com/golang/go.git
cd ../..
docker build -t golang:master --no-cache docker/golang-git
```

## Install Go from source package

Download, unpack and build wanted version using existing Go installation.
This example uses a local go 1.16.6 when bootstrapping.
```
wget https://go.dev/dl/go1.17.8.src.tar.gz
sudo tar xzf go1.17.8.src.tar.gz -C /usr/local/
cd /usr/local/go/src
sudo GOROOT_BOOTSTRAP=$HOME/sdk/go1.16.6 ./make.bash
```

## Uninstall

`sudo rm -rf /usr/local/go` and remove from path
