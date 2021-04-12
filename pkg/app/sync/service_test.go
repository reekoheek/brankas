package sync

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestService_Push(t *testing.T) {
	table := []struct {
		name string
		aDTO PushDTO
		xErr string
	}{
		// {
		// 	"positive 1",
		// 	PushDTO{},
		// 	"",
		// },
		// {
		// 	"positive 1",
		// 	PushDTO{},
		// 	"",
		// },
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{}

			err := s.Push(context.TODO(), tt.aDTO)
			if err != nil {
				assert.Equal(t, tt.xErr, err.Error())
				return
			}

			assert.Equal(t, "", tt.xErr, "expected err %s", tt.xErr)
		})
	}
}
