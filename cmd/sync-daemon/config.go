package main

import (
	"flag"
	"os"

	"github.com/xdlinux/sync.go/logs"
	"gopkg.in/yaml.v3"
)

var config struct {
	Daemon struct {
		TZ string `yaml:"timezone"`
	} `yaml:"daemon"`
}

func init() {
	parseFlags.Do(flag.Parse)
	data, err := os.ReadFile(*conf)
	if err != nil {
		logs.Error("Error reading config file, read error", map[string]string{
			"error": err.Error(),
		})
		panic("Error reading config file")
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logs.Error("Error reading config file, unmarshal error", map[string]string{
			"error": err.Error(),
		})
		panic("Error reading config file")
	}
}
