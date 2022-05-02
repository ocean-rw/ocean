package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Load(cfg interface{}, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}
