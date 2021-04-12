package acl

import (
	"testing"

	"gotest.tools/assert"
)

func TestACL_WriteBy(t *testing.T) {
	theACL := ACL{
		modes: map[string]int{
			"foo": READWRITE,
			"bar": READONLY,
			"baz": NONE,
		},
	}

	table := []struct {
		name      string
		aUsername string
		xAllowed  bool
	}{
		{
			"if explicitly rw",
			"foo",
			true,
		},
		{
			"if explicitly r",
			"bar",
			false,
		},
		{
			"if explicitly none",
			"baz",
			false,
		},
		{
			"if unspecified",
			"unknown",
			false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.xAllowed, theACL.WriteBy(tt.aUsername))
		})
	}
}

func TestACL_ReadBy(t *testing.T) {
	theACL := ACL{
		modes: map[string]int{
			"foo": READWRITE,
			"bar": READONLY,
			"baz": NONE,
		},
	}

	table := []struct {
		name      string
		aUsername string
		xAllowed  bool
	}{
		{
			"if explicitly rw",
			"foo",
			true,
		},
		{
			"if explicitly r",
			"bar",
			true,
		},
		{
			"if explicitly unknown",
			"baz",
			false,
		},
		{
			"if unspecified",
			"unknown",
			false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.xAllowed, theACL.ReadBy(tt.aUsername))
		})
	}
}

func TestACL_Put(t *testing.T) {
	table := []struct {
		name  string
		aMode int
		xErr  string
	}{
		{
			"set none",
			NONE,
			"",
		},
		{
			"set ro",
			READONLY,
			"",
		},
		{
			"set rw",
			READWRITE,
			"",
		},
		{
			"set unknown",
			100,
			"invalid mode",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			acl := ACL{
				modes: map[string]int{},
			}

			err := acl.Put("foo", tt.aMode)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
			assert.Equal(t, tt.aMode, acl.modes["foo"])
		})
	}
}
