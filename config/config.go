package config

import (
	"os"

	"github.com/xdlinux/sync.go/logs"
	"gopkg.in/yaml.v3"
)

var readFile = os.ReadFile // for mocking

type Config struct {
	Daemon struct {
		TZ     string `yaml:"timezone"`
		Secret string `yaml:"secret"`
	} `yaml:"daemon"`
	Jobs *JobConfig
}

func Parse(file string) *Config {
	var config Config
	data, err := readFile(file)
	if err != nil {
		logs.Error("Error reading config file, read error", map[string]string{
			"error": err.Error(),
			"file":  file,
		})
		panic("Error reading config file")
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logs.Error("Error reading config file, unmarshal error", map[string]string{
			"error": err.Error(),
			"file":  file,
		})
		panic("Error reading config file")
	}
	return &config
}
