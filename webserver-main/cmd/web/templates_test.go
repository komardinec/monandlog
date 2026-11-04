package main

import (
	"testing"
	"time"
)

func TestNiceDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 9, 10, 12, 23, 0, 0, time.UTC),
			want: "10 Sep 2025 at 12:23",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2025, 9, 10, 12, 23, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "10 Sep 2025 at 11:23",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd := NiceDate(tt.tm)

			if nd != tt.want {
				t.Errorf("got %q; but, want %q", nd, tt.want)
			}
		})
	}
}
