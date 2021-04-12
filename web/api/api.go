package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type api struct {
}

func (a *api) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// r.Post("/sync/put", a.hSyncPut())
	// r.Post("/sync/get", a.hSyncGet())
	// r.Get("/sync/nodes", h.Wrap(api.hUserNodes()))

	return r
}
