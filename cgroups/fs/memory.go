package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type memoryGroup struct {
}

func (s *memoryGroup) Stats(d *data) (map[string]float64, error) {
	paramData := make(map[string]float64)
	path, err := d.path("memory")
	if err != nil {
		return nil, err
	}

	// Set stats from memory.stat.
	statsFile, err := os.Open(filepath.Join(path, "memory.stat"))
	if err != nil {
		return nil, err
	}
	defer statsFile.Close()

	sc := bufio.NewScanner(statsFile)
	for sc.Scan() {
		t, v, err := getCgroupParamKeyValue(sc.Text())
		if err != nil {
			return nil, err
		}
		paramData[t] = v
	}

	// Set memory usage and max historical usage.
	params := []string{
		"usage_in_bytes",
		"max_usage_in_bytes",
	}
	for _, param := range params {
		value, err := getCgroupParamFloat64(path, fmt.Sprintf("memory.%s", param))
		if err != nil {
			return nil, err
		}
		paramData[param] = value
	}

	return paramData, nil
}
