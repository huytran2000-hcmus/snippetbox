package models

import (
	"testing"

	"github.com/huytran2000-hcmus/snippetbox/internal/assert"
)

func TestUserExists(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want bool
	}{
		{
			name: "Valid ID",
			id:   1,
			want: true,
		},
		{
			name: "Zero ID",
			id:   0,
			want: false,
		},
		{
			name: "Non-existent ID",
			id:   2,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := &UserDB{db}
			got, _ := m.Exists(tt.id)

			assert.Equal(t, got, tt.want)
		})
	}
}
