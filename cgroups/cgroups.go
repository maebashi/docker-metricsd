package cgroups

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/dotcloud/docker/pkg/mount"
)

var (
	ErrNotFound = errors.New("mountpoint not found")
)

type Cgroup struct {
	Name              string `json:"name,omitempty"`
	Parent            string `json:"parent,omitempty"`
	DeviceAccess      bool   `json:"device_access,omitempty"`      // name of parent cgroup or slice
	Memory            int64  `json:"memory,omitempty"`             // Memory limit (in bytes)
	MemoryReservation int64  `json:"memory_reservation,omitempty"` // Memory reservation or soft_limit (in bytes)
	MemorySwap        int64  `json:"memory_swap,omitempty"`        // Total memory usage (memory + swap); set `-1' to disable swap
	CpuShares         int64  `json:"cpu_shares,omitempty"`         // CPU shares (relative weight vs. other containers)
	CpuQuota          int64  `json:"cpu_quota,omitempty"`          // CPU hardcap limit (in usecs). Allowed cpu time in a given period.
	CpuPeriod         int64  `json:"cpu_period,omitempty"`         // CPU period to be used for hardcapping (in usecs). 0 to use system default.
	CpusetCpus        string `json:"cpuset_cpus,omitempty"`        // CPU to use
	Freezer           string `json:"freezer,omitempty"`            // set the freeze value for the process

	Slice string `json:"slice,omitempty"` // Parent slice to use for systemd
}

// https://www.kernel.org/doc/Documentation/cgroups/cgroups.txt
func FindCgroupMountpoint(subsystem string) (string, error) {
	mounts, err := mount.GetMounts()
	if err != nil {
		return "", err
	}

	for _, mount := range mounts {
		if mount.Fstype == "cgroup" {
			for _, opt := range strings.Split(mount.VfsOpts, ",") {
				if opt == subsystem {
					return mount.Mountpoint, nil
				}
			}
		}
	}
	return "", ErrNotFound
}

func GetInitCgroupDir(subsystem string) (string, error) {
	f, err := os.Open("/proc/1/cgroup")
	if err != nil {
		return "", err
	}
	defer f.Close()

	return parseCgroupFile(subsystem, f)
}

func parseCgroupFile(subsystem string, r io.Reader) (string, error) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}
		text := s.Text()
		parts := strings.Split(text, ":")
		for _, subs := range strings.Split(parts[1], ",") {
			if subs == subsystem {
				return parts[2], nil
			}
		}
	}
	return "", ErrNotFound
}
