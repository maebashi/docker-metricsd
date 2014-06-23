package fs

import (
	"bufio"
	"os"
	"path/filepath"

	"../../cgroups"
)

type memoryGroup struct {
}

func (s *memoryGroup) GetStats(d *data, stats *cgroups.Stats) error {
	path, err := d.path("memory")
	if err != nil {
		return err
	}

	// Set stats from memory.stat.
	statsFile, err := os.Open(filepath.Join(path, "memory.stat"))
	if err != nil {
		return err
	}
	defer statsFile.Close()

	sc := bufio.NewScanner(statsFile)
	for sc.Scan() {
		t, v, err := getCgroupParamKeyValue(sc.Text())
		if err != nil {
			return err
		}
		stats.MemoryStats.Stats[t] = v
	}

	// Set memory usage and max historical usage.
	value, err := getCgroupParamInt(path, "memory.usage_in_bytes")
	if err != nil {
		return err
	}
	stats.MemoryStats.Usage = value
	value, err = getCgroupParamInt(path, "memory.max_usage_in_bytes")
	if err != nil {
		return err
	}
	stats.MemoryStats.MaxUsage = value
	value, err = getCgroupParamInt(path, "memory.failcnt")
	if err != nil {
		return err
	}
	stats.MemoryStats.Failcnt = value

	return nil
}
