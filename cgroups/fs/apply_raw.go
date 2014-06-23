package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"../../cgroups"
)

var (
	subsystems = map[string]subsystem{
		"memory":  &memoryGroup{},
		"cpu":     &cpuGroup{},
		"cpuacct": &cpuacctGroup{},
		"blkio":   &blkioGroup{},
		"freezer": &freezerGroup{},
	}
)

type subsystem interface {
	//Set(*data) error
	//Remove(*data) error
	GetStats(*data, *cgroups.Stats) error
}

type data struct {
	root   string
	cgroup string
	c      *cgroups.Cgroup
	pid    int
}

func GetStats(c *cgroups.Cgroup) (*cgroups.Stats, error) {
	stats := cgroups.NewStats()

	d, err := getCgroupData(c, 0)
	if err != nil {
		return nil, err
	}

	for _, sys := range subsystems {
		if err := sys.GetStats(d, stats); err != nil {
			//return nil, err
			continue
		}
	}

	return stats, nil
}

func getCgroupData(c *cgroups.Cgroup, pid int) (*data, error) {
	// we can pick any subsystem to find the root
	cgroupRoot, err := cgroups.FindCgroupMountpoint("cpu")
	if err != nil {
		return nil, err
	}
	cgroupRoot = filepath.Dir(cgroupRoot)

	if _, err := os.Stat(cgroupRoot); err != nil {
		return nil, fmt.Errorf("cgroups fs not found")
	}

	cgroup := c.Name
	if c.Parent != "" {
		cgroup = filepath.Join(c.Parent, cgroup)
	}

	return &data{
		root:   cgroupRoot,
		cgroup: cgroup,
		c:      c,
		pid:    pid,
	}, nil
}

func (raw *data) parent(subsystem string) (string, error) {
	initPath, err := cgroups.GetInitCgroupDir(subsystem)
	if err != nil {
		return "", err
	}
	return filepath.Join(raw.root, subsystem, initPath), nil
}

func (raw *data) path(subsystem string) (string, error) {
	parent, err := raw.parent(subsystem)
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, raw.cgroup), nil
}

func (raw *data) join(subsystem string) (string, error) {
	path, err := raw.path(subsystem)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(path, 0755); err != nil && !os.IsExist(err) {
		return "", err
	}
	if err := writeFile(path, "cgroup.procs", strconv.Itoa(raw.pid)); err != nil {
		return "", err
	}
	return path, nil
}

func writeFile(dir, file, data string) error {
	return ioutil.WriteFile(filepath.Join(dir, file), []byte(data), 0700)
}
