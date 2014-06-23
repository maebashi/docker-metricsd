package fs

import (
	"bufio"
	"os"
	"path/filepath"
	"syscall"

	"../../cgroups"
)

type cpuGroup struct {
}

func (s *cpuGroup) GetStats(d *data, stats *cgroups.Stats) error {
	path, err := d.path("cpu")
	if err != nil {
		return err
	}

	f, err := os.Open(filepath.Join(path, "cpu.stat"))
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ENOENT {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		t, v, err := getCgroupParamKeyValue(sc.Text())
		if err != nil {
			return err
		}
		switch t {
		case "nr_periods":
			stats.CpuStats.ThrottlingData.Periods = v

		case "nr_throttled":
			stats.CpuStats.ThrottlingData.ThrottledPeriods = v

		case "throttled_time":
			stats.CpuStats.ThrottlingData.ThrottledTime = v
		}
	}
	return nil
}
