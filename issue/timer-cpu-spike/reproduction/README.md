# Reproduction of timer CPU issue

```
# Create certs for TLS in current directory
./gen-test-certs.sh

# Build server and client
go build server.go
go build client.go

# Start server and client (hardcoded certs paths)
./server
./client

# Get PID of the client process
ps -e | grep client

# Observe the CPU on the client process
pidstat -p <PID> 1

# The issue can often be seen after 30s to 7min

11:45:26 AM   UID       PID    %usr %system  %guest   %wait    %CPU   CPU  Command
11:45:27 AM  1000   2713246    0,00    0,00    0,00    0,00    0,00     2  client
11:45:28 AM  1000   2713246    0,00    0,00    0,00    0,00    0,00     2  client
11:45:29 AM  1000   2713246    0,00    0,00    0,00    0,00    0,00     2  client
11:45:30 AM  1000   2713246   56,00   12,00    0,00    0,00   68,00     2  client
11:45:31 AM  1000   2713246   81,00   19,00    0,00    0,00  100,00     2  client
11:45:32 AM  1000   2713246   78,00   23,00    0,00    0,00  101,00     2  client
11:45:33 AM  1000   2713246   78,00   22,00    0,00    0,00  100,00     2  client
11:45:34 AM  1000   2713246   77,00   22,00    0,00    0,00   99,00     2  client
11:45:35 AM  1000   2713246   24,00    6,00    0,00    0,00   30,00     2  client
11:45:36 AM  1000   2713246    0,00    0,00    0,00    0,00    0,00     2  client
11:45:37 AM  1000   2713246    1,00    0,00    0,00    0,00    1,00     2  client
11:45:38 AM  1000   2713246    0,00    0,00    0,00    0,00    0,00     2  client
```
