package config

import (
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/xdlinux/sync.go/logs"
	"gopkg.in/yaml.v3"
)

type JobConfig struct {
	ConfigDir  string `yaml:"config_dir,omitempty"`
	Cooldown   int    `yaml:"cooldown,omitempty"`
	HelpPath   string `yaml:"help,omitempty"`
	MirrorPath string `yaml:"url,omitempty"`
}

func (v *JobConfig) merge(t *JobConfig) {
	rv := reflect.ValueOf(v).Elem()
	rt := reflect.ValueOf(t).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).IsZero() {
			rv.Field(i).Set(rt.Field(i))
		}
	}
}

func (v *JobConfig) Format(job string) *JobConfig {
	rv := reflect.ValueOf(v).Elem()
	ret := &JobConfig{}
	rt := reflect.ValueOf(ret).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).Kind() == reflect.String {
			rt.Field(i).Set(reflect.ValueOf(
				strings.ReplaceAll(rv.Field(i).String(), "%i", job),
			))
		} else {
			rt.Field(i).Set(rv.Field(i))
		}
	}
	return ret
}

func (v *JobConfig) Get(job string) *JobConfig {
	// if value.Content
	path := path.Join(v.ConfigDir, job+".yml")
	file, err := readFile(path)
	if err != nil {
		logs.Warn("unable to find job config file, falling back to default", logs.T{
			"job": job, "error": err.Error(), "file": path,
		})
		return v.Format(job)
	}
	ret := &JobConfig{}
	err = yaml.Unmarshal(file, ret)
	if err != nil {
		logs.Error("Error reading config file, unmarshal error, falling back to default", map[string]string{
			"job": job, "error": err.Error(), "file": path,
		})
		return v.Format(job)
	}
	ret.merge(v)
	return ret.Format(job)
}

func (v *JobConfig) Save(job string) error {
	cfgdir := v.ConfigDir
	defer func() { v.ConfigDir = cfgdir }()
	v.ConfigDir = ""
	path := path.Join(cfgdir, job+".yml")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0644)
	if os.IsExist(err) {
		logs.Error("created config file must not exist", logs.T{
			"job": job, "error": err.Error(), "file": path,
		})
		return err
	}
	data, err := yaml.Marshal(v)
	if err != nil {
		logs.Error("yaml marshal error", logs.T{
			"job": job, "error": err.Error(), "file": path,
		})
		return err
	}
	for len(data) != 0 {
		n, err := file.Write(data)
		data = data[n:]
		if err != nil {
			logs.Error("file write error", logs.T{
				"job": job, "error": err.Error(), "file": path,
			})
		}
	}
	return nil
}
