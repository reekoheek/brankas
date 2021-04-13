package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reekoheek/brankas/pkg/app/sync"
	"gotest.tools/assert"
)

func TestAPI_hSyncPush(t *testing.T) {
	table := []struct {
		name       string
		dPusherErr error
		xStatus    int
		xErr       string
	}{
		{
			"positive case",
			nil,
			201,
			"",
		},
		{
			"if pusher err",
			fmt.Errorf("pusher err"),
			400,
			"pusher err",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			dto := sync.PushDTO{}
			a := &api{
				sPusher: tPusher(func(context.Context, sync.PushDTO) error {
					if tt.dPusherErr != nil {
						return tt.dPusherErr
					}
					return nil
				}),
			}
			h := a.hSyncPush()
			rr := httptest.NewRecorder()
			r := tNewRequest("POST", "/", dto)
			h(rr, r)
			assert.Equal(t, tt.xStatus, rr.Code)
			tAssertErr(t, tt.xErr, rr)
		})
	}
}

func TestAPI_hSyncPull(t *testing.T) {
	table := []struct {
		name          string
		dPullerResult []sync.EventDTO
		dPullerErr    error
		aDTO          interface{}
		xStatus       int
		xErr          string
	}{
		{
			"positive case 1",
			nil,
			nil,
			sync.PullDTO{},
			200,
			"",
		},
		{
			"positive case 2",
			[]sync.EventDTO{
				{Version: 1}, {Version: 2}, {Version: 3},
			},
			nil,
			sync.PullDTO{},
			200,
			"",
		},
		{
			"if puller err",
			nil,
			fmt.Errorf("puller err"),
			sync.PullDTO{},
			400,
			"puller err",
		},
		{
			"if empty body",
			nil,
			nil,
			nil,
			400,
			"invalid body",
		},
		{
			"if invalid body",
			nil,
			nil,
			"invalid",
			400,
			"invalid body",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			a := &api{
				sPuller: tPuller(func(context.Context, sync.PullDTO) ([]sync.EventDTO, error) {
					if tt.dPullerErr != nil {
						return nil, tt.dPullerErr
					}
					return tt.dPullerResult, nil
				}),
			}
			h := a.hSyncPull()
			rr := httptest.NewRecorder()
			r := tNewRequest("POST", "/", tt.aDTO)
			h(rr, r)
			assert.Equal(t, tt.xStatus, rr.Code)
			tAssertErr(t, tt.xErr, rr)
			if rr.Code >= 300 {
				return
			}

			result := []sync.EventDTO{}
			if err := parseBody(rr.Body, &result); err != nil {
				t.Error(err)
				return
			}

			assert.DeepEqual(t, tt.dPullerResult, result)
		})
	}
}

type tPusher func(ctx context.Context, dto sync.PushDTO) error

func (t tPusher) Push(ctx context.Context, dto sync.PushDTO) error {
	return t(ctx, dto)
}

type tPuller func(context.Context, sync.PullDTO) ([]sync.EventDTO, error)

func (t tPuller) Pull(ctx context.Context, dto sync.PullDTO) ([]sync.EventDTO, error) {
	return t(ctx, dto)
}

func tNewRequest(method string, url string, data interface{}) *http.Request {
	var body io.Reader
	if data != nil {
		bb, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(bb)
	}

	r, err := http.NewRequest(method, "/", body)
	if err != nil {
		panic(err)
	}

	return r
}

func tAssertErr(t *testing.T, xErr string, rr *httptest.ResponseRecorder) {
	if rr.Code < 300 {
		assert.Equal(t, "", xErr, "expected err %s", xErr)
		return
	}

	var httpErr httpError
	if err := parseBody(rr.Body, &httpErr); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, xErr, httpErr.Error())
}
