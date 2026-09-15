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
	case errors.Is(err, application.ErrForbidden):
		return http.StatusForbidden, http.StatusText(http.StatusForbidden)
	default:
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}
}
