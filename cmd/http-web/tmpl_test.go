package main

import (
	"testing"
	"time"
)

func TestFmtDate(t *testing.T) {

	tsts := []struct {
		name string
		tm   time.Time
		exp  string
	}{
		{
			name: "testUTC",
			tm:   time.Date(2023, 3, 4, 16, 5, 0, 0, time.UTC),
			exp:  "Mar 4 16:05 UTC 2023",
		},
		{
			name: "Emptytest",
			tm:   time.Time{},
			exp:  "",
		},
		{
			name: "testSGT",
			tm:   time.Date(2023, 1, 2, 15, 4, 5, 0, time.FixedZone("SGT", 8*60*60)),
			exp:  "Jan 2 07:04 UTC 2023",
		},
	}

	for _, tst := range tsts {
		t.Run(tst.name, func(t *testing.T) {
			fd := fmtDate(tst.tm)
			if fd != tst.exp {
				t.Errorf("Expected %s, got %s", tst.exp, fd)
			}
		})
	}

}
