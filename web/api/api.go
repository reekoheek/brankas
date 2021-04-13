package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/reekoheek/brankas/pkg/app/sync"
)

type api struct {
	sPusher sync.Pusher
	sPuller sync.Puller
}

func New(ss sync.Service) *api {
	return &api{
		sPusher: ss,
		sPuller: ss,
	}
}

func (a *api) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Post("/sync/push", a.hSyncPush())
	r.Post("/sync/pull", a.hSyncPull())
	// r.Get("/sync/nodes", h.Wrap(api.hUserNodes()))

	return r
}
