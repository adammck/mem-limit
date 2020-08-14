# Mem Limit

This is a quick experiment to determine how much memory is available when
running inside a cgroup, or (optionally) falling back to the host memory. I
might turn it into a library later.

## Usage

```console
$ go run main.go
[info] not running in a container
available host memory: 17179869184 bytes

$ docker build -t mem-limit .
[...]

$ docker run mem-limit:latest
[warn] in container but no memory limit set
available host memory: 2087837696 bytes

$ docker run --memory=134217728 mem-limit:latest
available container memory: 134217728 bytes
```

## License

MIT.
