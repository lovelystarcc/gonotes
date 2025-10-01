package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewErrResponse(status int, err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: status,
		Message:        err.Error(),
	}
}
