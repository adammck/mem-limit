package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	// TODO(adammck): Remove this dep, only using a single attribute.
	"github.com/shirou/gopsutil/mem"
)

func hostMemoryLimit() (uint64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return vm.Total, nil
}

var errNotInContainer = errors.New("not running in a container")
var errNoMemoryLimit = errors.New("no cgroup memory limit is set")

func containerMemoryLimit() (uint64, error) {
	fn := "/sys/fs/cgroup/memory/memory.limit_in_bytes"

	// return early if the file doesn't exist. we're probably just not running
	// in a cgroup.
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return 0, errNotInContainer
	}

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseUint(strings.TrimSuffix(string(buf), "\n"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("while casting cgroup memory limit to uint64: %s", err)
	}

	// return an error instead of the special *no limit* value. callers will
	// probably want to fall back to the host memory limit rather than assume
	// they have 8EiB of memory available.
	if val == 9223372036854771712 {
		return val, errNoMemoryLimit
	}

	return val, nil
}

func main() {
	lim, err := containerMemoryLimit()
	if err == nil {
		fmt.Printf("available container memory: %d bytes\n", lim)
		os.Exit(0)
	}

	if err == errNotInContainer {
		fmt.Println("[info] not running in a container")

	} else if err == errNoMemoryLimit {
		fmt.Println("[warn] in container but no memory limit set")

	} else {
		fmt.Printf("error fetching container memory: %s\n", err)
		os.Exit(1)
	}

	lim, err = hostMemoryLimit()
	if err != nil {
		fmt.Printf("error fetching host memory: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("available host memory: %d bytes\n", lim)
	os.Exit(0)
}
