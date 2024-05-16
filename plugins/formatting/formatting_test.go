package formatting

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	d0 := time.Duration(43) * time.Second
	d0s := "43s"
	if Duration(d0) != d0s {
		t.Errorf("Duration(%v) = %s; want %s", d0, Duration(d0), d0s)
	}
	d1 := time.Duration(95) * time.Second
	d1s := "1m 35s"
	if Duration(d1) != d1s {
		t.Errorf("Duration(%v) = %s; want %s", d1, Duration(d1), d1s)
	}
	d2 := time.Duration(3855) * time.Second
	d2s := "1h 4m"
	if Duration(d2) != d2s {
		t.Errorf("Duration(%v) = %s; want %s", d2, Duration(d2), d2s)
	}
	d3 := time.Duration(589442) * time.Second
	d3s := "6d 19h"
	if Duration(d3) != d3s {
		t.Errorf("Duration(%v) = %s; want %s", d3, Duration(d3), d3s)
	}
	d4 := time.Duration(4589442) * time.Second
	d4s := "53d"
	if Duration(d4) != d4s {
		t.Errorf("Duration(%v) = %s; want %s", d4, Duration(d4), d4s)
	}
}
