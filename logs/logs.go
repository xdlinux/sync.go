package logs

import (
	"bytes"
	"fmt"
	"log"

	"github.com/coreos/go-systemd/v22/journal"
)

var dispatch func(msg string, level journal.Priority, tags map[string]string) error

type T = map[string]string

func formatTags(tags map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range tags {
		fmt.Fprintf(b, " %s=%s", key, value)
	}
	return b.String()
}

func init() {
	if journal.Enabled() {
		dispatch = journal.Send
	} else {
		log.Printf("[WARN] msg=journalctl not enabled, falling back to stdout/stderr")
		dispatch = func(msg string, level journal.Priority, tags map[string]string) error {
			log.Printf("[%s]%s msg=%s", map[journal.Priority]string{
				3: "ERROR", 4: "WARN", 6: "INFO",
			}[level], formatTags(tags), msg)
			return nil
		}
	}
}

func ensureDispatch(msg string, level journal.Priority, tags map[string]string) {
	var err error
	for ok, i := false, 10; !ok && i > 0; ok, i = err == nil, i-1 {
		err = dispatch(msg, level, tags)
	}
}

func Info(msg string, tags map[string]string) {
	go ensureDispatch(msg, journal.PriInfo, tags)
}

func Warn(msg string, tags map[string]string) {
	go ensureDispatch(msg, journal.PriWarning, tags)
}

func Error(msg string, tags map[string]string) {
	go ensureDispatch(msg, journal.PriErr, tags)
}
