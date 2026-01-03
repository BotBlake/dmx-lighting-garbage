package device

import (
	"encoding/xml"
	"os"
	"path/filepath"
)

func LoadProfiles(dir string) (map[string]*DDFDevice, error) {
	profiles := make(map[string]*DDFDevice)

	files, err := filepath.Glob(filepath.Join(dir, "*.xml"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var ddf DDFDevice
		if err := xml.Unmarshal(data, &ddf); err != nil {
			continue
		}

		// Use filename (without extension) as key
		name := filepath.Base(file)
		name = name[:len(name)-len(filepath.Ext(name))]

		profiles[name] = &ddf
	}

	return profiles, nil
}
