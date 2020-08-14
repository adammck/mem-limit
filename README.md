# Mem Limit

This is a quick experiment to determine how much memory is available and in use
when running inside a cgroup, or (optionally) falling back to the host memory. I
might turn it into a library later.

## Usage

```console
$ go run main.go
[info] not running in a container
total: 17179869184 bytes
used: 11055026176 bytes

$ docker build -t mem-limit .
[...]

$ docker run mem-limit:latest
[warn] in container but no memory limit set; falling back to host metrics
total: 2087837696 bytes
used: 347717632 bytes

$ docker run --memory 128M mem-limit:latest
total: 134217728 bytes
used: 1466368 bytes
```

## License

MIT.
