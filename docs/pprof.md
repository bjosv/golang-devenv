# pprof - visualization and analysis of profiling data

## Use pprof towards K8s

```
kubectl port-forward <pod> 6060 &
go tool pprof http://localhost:6060/debug/pprof/profile
go tool pprof -http=":8080" <generated-file-path>
```

## Use perf and pprof to view results

First build and install in path:
https://github.com/google/perf_data_converter

```
# Profile running process
ps -efa | grep <name>
sudo perf record -F 999 -g -p <pid>
sudo chmod 755 perf.data

# Get binary with debug symbols
kubectl cp <pod>:/redis_exporter redis_exporter -c redis-exporter

# Visualise result
go tool pprof -http=":8080" redis_exporter perf.data
```

## Trace

```
# Add
import "github.com/pkg/profile"
defer profile.Start(profile.TraceProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()

# Inspect results
go tool trace trace.out
```
