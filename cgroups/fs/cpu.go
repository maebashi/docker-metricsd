package fs

import (
	"bufio"
	"os"
	"path/filepath"
)

type cpuGroup struct {
}

func (s *cpuGroup) Stats(d *data) (map[string]float64, error) {
	paramData := make(map[string]float64)
	path, err := d.path("cpu")
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(path, "cpu.stat"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		t, v, err := getCgroupParamKeyValue(sc.Text())
		if err != nil {
			return nil, err
		}
		paramData[t] = v
	}
	return paramData, nil
}
