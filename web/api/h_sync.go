package api

import (
	"net/http"

	"github.com/reekoheek/brankas/pkg/app/sync"
	"github.com/reekoheek/brankas/web"
)

func (a *api) hSyncPush() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto sync.PushDTO

		if err := web.ParseBody(r.Body, &dto); err != nil {
			web.RespondErr(w, 400, err)
			return
		}

		if err := a.sPusher.Push(r.Context(), dto); err != nil {
			web.RespondErr(w, 400, err)
			return
		}

		web.Respond(w, 201, nil)
	}
}

func (a *api) hSyncPull() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto sync.PullDTO

		if err := web.ParseBody(r.Body, &dto); err != nil {
			web.RespondErr(w, 400, err)
			return
		}

		result, err := a.sPuller.Pull(r.Context(), dto)
		if err != nil {
			web.RespondErr(w, 400, err)
			return
		}

		web.Respond(w, 200, result)
	}
}
