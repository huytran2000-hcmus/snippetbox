package main

import (
	"testing"
	"time"

	"github.com/huytran2000-hcmus/snippetbox/internal/assert"
)

func TestReadableDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "Thursday, 17 Mar 2022 at 10:15:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "Thursday, 17 Mar 2022 at 09:15:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := readableDate(tt.tm)
			assert.Equal(t, rd, tt.want)
		})
	}
}
