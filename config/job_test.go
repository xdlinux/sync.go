package config

import (
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func init() {
	fs := fstest.MapFS{
		"conf.d/test.yml": {
			Data: []byte(`cooldown: 20`),
		},
	}
	readFile = fs.ReadFile
}

func TestJobConfigGet(t *testing.T) {
	c1 := &JobConfig{
		ConfigDir:  "./conf.d",
		Cooldown:   10,
		HelpPath:   "%i_test/asdf",
		MirrorPath: "asdfa%i",
	}
	want := &JobConfig{
		ConfigDir:  "./conf.d", // unnecessary, but don't bother to reset it
		Cooldown:   20,
		HelpPath:   "test_test/asdf",
		MirrorPath: "asdfatest",
	}
	if diff := cmp.Diff(want, c1.Get("test")); diff != "" {
		t.Errorf("JobConfigGet mismatch (-want +got):\n%s", diff)
	}
}

func TestJobConfigFormat(t *testing.T) {
	c1 := &JobConfig{
		Cooldown:   10,
		HelpPath:   "%i_test/asdf",
		MirrorPath: "asdfa%i",
	}
	want := &JobConfig{
		Cooldown:   10,
		HelpPath:   "test_test/asdf",
		MirrorPath: "asdfatest",
	}
	if diff := cmp.Diff(want, c1.Format("test")); diff != "" {
		t.Errorf("JobConfigFormat mismatch (-want +got):\n%s", diff)
	}
}
