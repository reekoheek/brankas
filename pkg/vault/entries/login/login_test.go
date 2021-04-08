package login

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogin_Expired(t *testing.T) {
	table := []struct {
		name     string
		dExpiry  time.Time
		xExpired bool
	}{
		{
			"positive 1",
			time.Now(),
			true,
		},
		{
			"positive 2",
			time.Now().Add(-1 * time.Second),
			true,
		},
		{
			"positive 3",
			time.Now().Add(1 * time.Second),
			false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			e := Login{Expiry: tt.dExpiry}
			assert.Equal(t, tt.xExpired, e.Expired())
		})
	}
}
