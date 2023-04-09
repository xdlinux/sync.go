package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"flag"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/xdlinux/sync.go/config"
	"github.com/xdlinux/sync.go/logs"
)

var config_path = flag.String("conf", "config.yml", "config file for sync-daemon")
var cfg *config.Config

func init() {
	flag.Parse()
	cfg = config.Parse(*config_path)
}

func main() {
	logs.Info("starting sync-daemon", nil)
	location, err := time.LoadLocation(cfg.Daemon.TZ)
	if err != nil {
		logs.Error("Unable to set location", nil)
		os.Exit(1)
	}
	time.Local = location

	listeners, _ := activation.Listeners()
	if len(listeners) == 0 {
		addr, _ := net.ResolveUnixAddr("unix", "./daemon.sock")
		l, _ := net.ListenUnix("unix", addr)
		listeners = append(listeners, l)
	}
	if len(listeners) != 1 {
		logs.Error("Unexpected number of socket activation fds", logs.T{
			"listeners": fmt.Sprintf("%d", len(listeners)),
		})
		os.Exit(1)
	}
	// TODO
}
