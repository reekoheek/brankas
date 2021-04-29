package auth

import (
	"net/http"

	"github.com/reekoheek/brankas/web"
)

func (a *auth) hLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meta := struct {
			Name string `json:"name"`
		}{
			Name: "brankas",
		}

		web.Respond(w, 202, meta)
	}
}
