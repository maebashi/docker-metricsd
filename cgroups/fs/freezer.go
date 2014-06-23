package fs

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"../../cgroups"
)

type freezerGroup struct {
}

func getFreezerFileData(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	return strings.TrimSuffix(string(data), "\n"), err
}

func (s *freezerGroup) GetStats(d *data, stats *cgroups.Stats) error {
	path, err := d.path("freezer")
	if err != nil {
		return err
	}
	var data string
	if data, err = getFreezerFileData(filepath.Join(path, "freezer.parent_freezing")); err != nil {
		return err
	}
	stats.FreezerStats.ParentState = data
	if data, err = getFreezerFileData(filepath.Join(path, "freezer.self_freezing")); err != nil {
		return err
	}
	stats.FreezerStats.SelfState = data

	return nil
}
