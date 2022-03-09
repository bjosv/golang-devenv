# Go TLS test

## Build test container
docker build -t bjosv/tls-issue-tester:1.16.7 -f Dockerfile.go1.16.7 .

minikube image load bjosv/tls-issue-tester:1.16.7
