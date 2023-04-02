package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"flag"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/xdlinux/sync.go/logs"
)

var conf = flag.String("conf", "config.yml", "config file for sync-daemon")
var parseFlags = sync.Once{}

func init() {
	parseFlags.Do(flag.Parse)
}

func main() {
	logs.Info("starting sync-daemon", nil)
	location, err := time.LoadLocation(config.Daemon.TZ)
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
