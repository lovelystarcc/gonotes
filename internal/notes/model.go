package notes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func (n *Note) Bind(r *http.Request) error {
	if n.Text == "" {
		return fmt.Errorf("text is required")
	}
	return nil
}

func (n *Note) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

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
