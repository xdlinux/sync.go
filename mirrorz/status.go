package mirrorz

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/xdlinux/sync.go/mirrorz/consts"
)

type MirrorStatus map[rune]time.Time

func (st MirrorStatus) mainState(s rune, t time.Time) {
	for _, s := range "SDYFP" {
		delete(st, s)
	}
	st[s] = t
}

func (st MirrorStatus) SuccessSince(since time.Time) {
	st.mainState('S', since)
}
func (st MirrorStatus) PendingSince(since time.Time) {
	st.mainState('D', since)
}
func (st MirrorStatus) SyncingSince(since time.Time) {
	st.mainState('Y', since)
}
func (st MirrorStatus) FailedSince(since time.Time) {
	st.mainState('F', since)
}
func (st MirrorStatus) PausedSince(since time.Time) {
	st.mainState('P', since)
}

func (st MirrorStatus) Flag(flag consts.MirrorStatusFlag) {
	st[rune(flag)] = time.Time{}
}

func (st MirrorStatus) NextScheduled(tm time.Time) {
	st['X'] = tm
}
func (st MirrorStatus) NewMirror(tm time.Time) {
	st['N'] = tm
}
func (st MirrorStatus) LastSuccess(tm time.Time) {
	st['O'] = tm
}

func (st MirrorStatus) UnmarshalText(txt []byte) error {
	stream := strings.NewReader(string(txt))
	for {
		ch, _, err := stream.ReadRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if ch <= 'A' || ch >= 'Z' {
			return fmt.Errorf("invalid status identifier %c in status string %s", ch, string(txt))
		}
		var ts strings.Builder
		for {
			ch, _, err := stream.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if ch >= '0' && ch <= '9' {
				ts.WriteRune(ch)
			} else {
				stream.UnreadRune()
				break
			}
		}
		if ts.Len() == 0 {
			st[ch] = time.Time{}
		} else {
			t, _ := strconv.ParseInt(ts.String(), 10, 64)
			st[ch] = time.Unix(t, 0)
		}
	}
}

func (st MirrorStatus) MarshalText() ([]byte, error) {
	var builder bytes.Buffer
	for _, s := range "SDYFP" {
		if v, ok := st[s]; ok && !v.IsZero() {
			builder.WriteRune(s)
			builder.WriteString(strconv.FormatInt(v.Unix(), 10))
			break
		}
	}
	for k, v := range st {
		if strings.ContainsRune("SDYFP", k) {
			continue
		}
		builder.WriteRune(k)
		if !v.IsZero() {
			builder.WriteString(strconv.FormatInt(v.Unix(), 10))
		}
	}
	return builder.Bytes(), nil
}
