package handlersdto

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorDTO struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ErrorDTO) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Status)
	return nil
}

func NewErrorDTO(status int, err error) *ErrorDTO {
	if err == nil {
		return &ErrorDTO{
			Status:  status,
			Message: "",
		}
	}
	return &ErrorDTO{
		Status:  status,
		Message: err.Error(),
	}
}
