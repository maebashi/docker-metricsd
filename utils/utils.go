package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/docker/libcontainer/cgroups"
	"github.com/docker/libcontainer/cgroups/fs"
	"github.com/docker/libcontainer/system"
)

var devDir string = ""
var parentName string = ""

func GetCgroupStats(id string) (m *cgroups.Stats, err error) {
	if id, err = getLongID(id); err != nil {
		return
	}
	c := cgroups.Cgroup{
		Parent: parentName,
		Name:   id,
	}
	return fs.GetStats(&c)
}

func getLongID(shortID string) (longID string, err error) {
	if devDir == "" {
		devDir, err = cgroups.FindCgroupMountpoint("devices")
		if err != nil {
			return
		}
	}

	pat := filepath.Join(devDir, "*", shortID+"*", "tasks")
	a, err := filepath.Glob(pat)
	if err != nil {
		return
	}
	if len(a) != 1 {
		return "", fmt.Errorf("No such ID %s", shortID)
	}
	dir := filepath.Dir(a[0])
	longID = filepath.Base(dir)
	parentName = filepath.Base(filepath.Dir(dir))
	return
}

func GetContainerPID(id string) (pid string, err error) {
	if devDir == "" {
		devDir, err = cgroups.FindCgroupMountpoint("devices")
		if err != nil {
			return
		}
	}

	pat := filepath.Join(devDir, "*", id+"*", "tasks")
	a, err := filepath.Glob(pat)
	if err != nil {
		return
	}
	if len(a) != 1 {
		return "", fmt.Errorf("No such ID %s", id)
	}

	contents, err := ioutil.ReadFile(a[0])
	if err != nil {
		return
	}

	a = strings.Split(string(contents), "\n")
	return a[0], nil
}

func NetNsSynchronize(pid string, fn func() error) (err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f, err := os.OpenFile("/proc/self/ns/net", os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer func() {
		system.Setns(f.Fd(), 0)
		f.Close()
	}()
	if err = setNetNs(pid); err != nil {
		return
	}
	return fn()
}

func setNetNs(pid string) (err error) {
	path := filepath.Join("/proc", pid, "ns", "net")
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()
	if err = system.Setns(f.Fd(), 0); err != nil {
		return
	}
	return
}

func GetIfStats() (m map[string]interface{}, err error) {
	m = map[string]interface{}{}
	cmd := exec.Command("cat", "/proc/net/dev")
	f, err := cmd.StdoutPipe()
	//f, err := os.Open("/proc/net/dev")
	if err != nil {
		return
	}
	defer f.Close()
	err = cmd.Start()
	if err != nil {
		return
	}

	s := bufio.NewScanner(f)
	var d uint64
	for i := 0; s.Scan(); {
		var name string
		var n [8]uint64
		text := s.Text()
		if strings.Index(text, ":") < 1 {
			continue
		}
		ts := strings.Split(text, ":")
		fmt.Sscanf(ts[0], "%s", &name)
		if strings.HasPrefix(name, "veth") || name == "lo" {
			continue
		}
		fmt.Sscanf(ts[1],
			"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
			&n[0], &n[1], &n[2], &n[3], &d, &d, &d, &d,
			&n[4], &n[5], &n[6], &n[7], &d, &d, &d, &d,
		)
		j := "." + strconv.Itoa(i)
		m["name"+j] = name
		m["inbytes"+j] = n[0]
		m["inpackets"+j] = n[1]
		m["inerrs"+j] = n[2]
		m["indrop"+j] = n[3]
		m["outbytes"+j] = n[4]
		m["outpackets"+j] = n[5]
		m["outerrs"+j] = n[6]
		m["outdrop"+j] = n[7]
		i++
	}
	err = cmd.Wait()
	return
}

func GetIfAddr(ifname string) net.Addr {
	i, err := net.InterfaceByName(ifname)
	if err == nil {
		addrs, err := i.Addrs()
		if err == nil {
			for _, addr := range addrs {
				if addr.(*net.IPNet).IP.DefaultMask() != nil {
					return addr
				}
			}
		}
	}
	return nil
}
