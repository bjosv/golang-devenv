# Timer cpu issue

These instructions setup a local Kubernetes cluster and pods using the tool `minikube`.

In the cluster there is pod `Prometheus` that every other minute connects to 9 Go processes
called `redis_exporter`, i.e the Go process that gives the CPU spike.
When `redis_exporter` accepts a http connection from `Prometheus` it connects to the redis database
via TLS to fetch metrics. This is the connection that sometimes triggers a "busy loop" in the Go runtime.

<prometheus>  ---->  <redis_exporter>  --X-->  <redis>

The official redis_exporter container used in this deployment is build with Go 1.17.8

## Requirements
- openssl  (used OpenSSL 1.1.1f)
- minikube (used v1.23.2, that gives K8s v1.22.2)
- kubectl

Tested on Ubuntu 20.04 using Docker 20.10.10

## Setup

### Generate TLS certificates

This script will generate TLS certs to `/tmp/tls-data` needed by the deployment.
The script is located relative to this readme-file in this git repo.

`./scripts/gen-test-certs.sh`

### Preparing a Kubernetes cluster

This will create a local K8s cluster with the TLS-config directory mounted.

```
minikube config set memory 8192
minikube config set cpus 6
minikube start --mount-string="/tmp/tls-data:/tls-data" --mount
kubectl get all -A
```

### Install prometheus on Kubernetes

This installs a standard tool that fetch metrics from installed pods.
This will trigger a http-connect to the Go program `redis_exporter` where we have the problem.

`kubectl create -f ./manifests/prometheus.yaml`

### Install redis_exporter (gives the Go issue) and redis in a pod

```
kubectl create -f manifests/redis-and-exporter-deployment.yaml
```

Wait for all pod status Running using following command:

`kubectl get pods`


## Wait for the CPU to spike

This step usually takes 5 to 30 minutes, sometimes more.

### Check the CPU on Go processes

Enter the VM/container that runs K8s:

`minikube ssh`

Watch for CPU spike using top

`docker@minikube:~$ top`

Example:
```
top - 12:41:07 up 1 day,  1:37,  0 users,  load average: 1.09, 0.87, 0.84
Tasks:  97 total,   1 running,  96 sleeping,   0 stopped,   0 zombie
%Cpu(s):  7.9 us,  2.6 sy,  0.0 ni, 85.6 id,  0.3 wa,  0.0 hi,  3.7 si,  0.0 st
MiB Mem :  31755.9 total,  16293.1 free,   5039.4 used,  10423.4 buff/cache
MiB Swap:    881.5 total,    881.5 free,      0.0 used.  25718.4 avail Mem

    PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
  49344 59000     20   0  712988  13936   7524 S 100.0   0.0   1:50.42 redis_exporter
   2312 root      20   0 2484280 168524  84968 S  10.7   0.5  17:09.58 kubelet
   1941 root      20   0 1105852 307484  75948 S   5.3   0.9  10:26.68 kube-apiserver
  49781 59000     20   0  713756  13884   7460 S   1.7   0.0   0:01.95 redis_exporter
...
```
