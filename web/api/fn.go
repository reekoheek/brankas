package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpError struct {
	Message  string       `json:"message"`
	Field    string       `json:"field,omitempty"`
	Children []*httpError `json:"children,omitempty"`
}

func (err httpError) Error() string {
	return err.Message
}

func (err *httpError) AddChild(child *httpError) {
	err.Children = append(err.Children, child)
}

func newError(err error) *httpError {
	switch v := err.(type) {
	case httpError:
		return &v
	case *httpError:
		return v
	default:
		return &httpError{
			Message: v.Error(),
		}
	}
}

func parseBody(body io.Reader, data interface{}) error {
	if body == nil {
		return fmt.Errorf("invalid body")
	}

	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("invalid body")
	}

	return nil
}

func respondErr(w http.ResponseWriter, statusCode int, err error) {
	if statusCode < 300 || statusCode >= 600 {
		statusCode = 500
	}

	respond(w, statusCode, newError(err))

}

func respond(w http.ResponseWriter, statusCode int, data interface{}) {
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
