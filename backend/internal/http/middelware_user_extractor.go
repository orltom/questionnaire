package http

import (
	"log/slog"
	"net/http"

	"gitlab.com/orltom/questionnaire/backend/internal/identity"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/infrastructure/rest"
)

func ExtractUser(next http.Handler, service *identity.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := service.Authenticate(ctx, "https://accounts.google.com", "1234", "demo")
		if err != nil {
			slog.ErrorContext(ctx, "failed to authenticate", "error", err)
			http.Error(w, "unexpected error", http.StatusInternalServerError)
			return
		}
		ctx = rest.WithActor(ctx, user.ID())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
