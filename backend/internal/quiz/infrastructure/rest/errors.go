package rest

import (
	"errors"
	"net/http"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
)

func httpError(err error) (int, string) {
	switch {
	case errors.Is(err, application.ErrInvalidArguments):
		return http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
	case errors.Is(err, application.ErrEntityNotFound):
		return http.StatusNotFound, http.StatusText(http.StatusNotFound)
	default:
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}
}
