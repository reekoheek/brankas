package api

import (
	"net/http"

	"github.com/reekoheek/brankas/web"
)

func (a *api) hIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		web.Respond(w, 202, struct {
			Name string `json:"name"`
		}{
			Name: "brankas",
		})
	}
}
