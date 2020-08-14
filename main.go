package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	// TODO(adammck): Remove this dep, only using two attributes.
	"github.com/shirou/gopsutil/mem"
)

type memInfo struct {
	Total uint64
	Used  uint64
}

func hostMemory() (memInfo, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return memInfo{}, err
	}

	return memInfo{
		Total: vm.Total,
		Used:  vm.Used,
	}, nil
}

var errNotInContainer = errors.New("not running in a container")
var errNoMemoryLimit = errors.New("no cgroup memory limit is set")

func containerMemory() (memInfo, error) {
	tot, err := readCgroupMemoryFile("limit_in_bytes")
	if err != nil {
		return memInfo{}, err
	}

	// return an error instead of the special *no limit* value. callers will
	// probably want to fall back to the host memory limit rather than assume
	// they have 8EiB of memory available.
	if tot == 9223372036854771712 {
		return memInfo{}, errNoMemoryLimit
	}

	use, err := readCgroupMemoryFile("usage_in_bytes")
	if err != nil {
		return memInfo{}, err
	}

	return memInfo{
		Total: tot,
		Used:  use,
	}, nil
}

func readCgroupMemoryFile(suffix string) (uint64, error) {
	fn := fmt.Sprintf("/sys/fs/cgroup/memory/memory.%s", suffix)

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
		return 0, fmt.Errorf("while casting cgroup memory to uint64: %s", err)
	}

	return val, nil
}

func main() {
	useHost := false
	mem, err := containerMemory()

	if err == errNotInContainer {
		fmt.Println("[info] not running in a container")
		useHost = true

	} else if err == errNoMemoryLimit {
		fmt.Println("[warn] in container but no memory limit set; falling back to host metrics")
		useHost = true

	} else if err != nil {
		fmt.Printf("error fetching container memory stats: %s\n", err)
		os.Exit(1)
	}

	if useHost {
		mem, err = hostMemory()
		if err != nil {
			fmt.Printf("error fetching host memory stats: %s\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("total: %d bytes\n", mem.Total)
	fmt.Printf("used: %d bytes\n", mem.Used)
}
