package auth

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/reekoheek/brankas/web"
)

type auth struct {
}

func New() *auth {
	return &auth{}
}

func (a *auth) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		web.RespondErr(w, 404, web.Error{Message: "not found"})
	})

	r.Post("/login", a.hLogin())

	return r
}
