# Delve - a debugger for Go

## Enabling delve in a container

In the Dockerfile for the debug target add:

```
# Add deps needed by delve
RUN apk update && apk upgrade && \
    apk add --no-cache git openssh build-base
...
# Build delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest
...
# Copy delve from the builder (when using multi-stage builds)
COPY --from=builder /go/bin/dlv /dlv
...
# Let delve execute the target binary (/redis_exporer) and start the debug session, but continue running.
CMD ["/dlv", "--listen=:40000", "--headless", "--api-version=2", "--accept-multiclient", "--log", "exec", "/redis_exporter", "--continue"]
```

Possibly make sure the root user is used.

## Building debug target

Use compiler flags [gcflags](https://pkg.go.dev/cmd/compile)
- Disable compiler optimizations and inlining: `-gcflags="all=-N -l"`

Use linker flags [ldflags](https://pkg.go.dev/cmd/link)
- Enable symbol tables: remove `-s -w` if used

## Update K8s deployment

Add SYS_PTRACE capability (and apparmor)
```
spec:
  template:
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/<CONTAINER-NAME>: unconfined
...
    spec:
      containers:
      - name: <CONTAINER-NAME>
        securityContext:
          capabilities:
            add:
            - SYS_PTRACE

```

## Attach to delve in container

```
kubectl port-forward <pod> 40001:40000 &
dlv connect :40001

# Change source path, example
config substitute-path /go/src/github.com/oliver006/ ~/git/
```

## Delve config

```
~/.config/dlv/config.yml
```
