package login

import (
	"testing"
	"time"

	"gotest.tools/assert"
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

func TestLogin_AddURL(t *testing.T) {
	table := []struct {
		name  string
		login Login
		aURL  string
		xErr  string
	}{
		{
			"positive 1",
			Login{},
			"http://foo.bar",
			"",
		},
		{
			"positive 2",
			Login{
				URLs: []string{"http://1.2"},
			},
			"http://foo.bar",
			"",
		},
		{
			"if url exists",
			Login{
				URLs: []string{"http://foo.bar"},
			},
			"http://foo.bar",
			"url exists",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			urlCount := len(tt.login.URLs)

			err := tt.login.AddURL(tt.aURL)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "Expected err %s", tt.xErr)
			assert.Equal(t, urlCount+1, len(tt.login.URLs))
		})
	}
}

func TestLogin_RemoveURL(t *testing.T) {
	table := []struct {
		name  string
		login Login
		aURL  string
		xErr  string
	}{
		{
			"positive 1",
			Login{
				URLs: []string{"http://foo.bar"},
			},
			"http://foo.bar",
			"",
		},
		{
			"if url not exists",
			Login{},
			"http://foo.bar",
			"url not exists",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			urlCount := len(tt.login.URLs)

			err := tt.login.RemoveURL(tt.aURL)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "Expected err %s", tt.xErr)
			assert.Equal(t, urlCount-1, len(tt.login.URLs))
		})
	}
}
