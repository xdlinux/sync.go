package mirrorz

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/xdlinux/sync.go/mirrorz/consts"
)

func (st MirrorStatus) MarshalTextStringBuilder() ([]byte, error) {
	var builder strings.Builder
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
	return []byte(builder.String()), nil
}

func BenchmarkMarshalText(b *testing.B) {
	st := MirrorStatus{
		'S': time.Now(),
		'C': time.Time{},
	}

	b.Run("original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := st.MarshalText(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := st.MarshalTextStringBuilder(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestUnmarshalStatus(t *testing.T) {
	tm := time.Now().Truncate(time.Second)
	str := fmt.Sprintf("\"S%dRO%d\"", tm.Unix(), tm.Unix())
	got := make(MirrorStatus)
	_ = json.Unmarshal([]byte(str), &got)
	want := MirrorStatus{}
	want.SuccessSince(tm)
	want.Flag(consts.ReverseProxy)
	want.LastSuccess(tm)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unmarshal mismatch (-want +got):\n%s", diff)
	}
}

func TestMarshalStatus(t *testing.T) {
	tm := time.Now().Truncate(time.Second)
	want := fmt.Sprintf("\"S%dRO%d\"", tm.Unix(), tm.Unix())
	tg := MirrorStatus{}
	tg.SuccessSince(tm)
	tg.Flag(consts.ReverseProxy)
	tg.LastSuccess(tm)
	got, err := json.Marshal(tg)
	if err != nil {
		t.Errorf("Marshal error: %v", err)
	}
	if diff := cmp.Diff([]byte(want), got); diff != "" {
		t.Errorf("Marshal mismatch (-want +got):\n%s", diff)
	}
}
