package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ParseBody(body io.Reader, data interface{}) error {
	if body == nil {
		return fmt.Errorf("invalid body")
	}

	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("invalid body")
	}

	return nil
}

func RespondErr(w http.ResponseWriter, statusCode int, err error) {
	if statusCode < 300 || statusCode >= 600 {
		statusCode = 500
	}

	Respond(w, statusCode, NewError(err))

}

func Respond(w http.ResponseWriter, statusCode int, data interface{}) {
	if statusCode < 200 || statusCode >= 600 {
		statusCode = 200
	}

	w.WriteHeader(statusCode)

	bb, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(bb); err != nil {
		panic(err)
	}
}
