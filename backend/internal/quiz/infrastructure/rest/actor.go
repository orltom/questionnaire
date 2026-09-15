package rest

import (
	"context"
	"errors"

	"gitlab.com/orltom/questionnaire/backend/internal/identity"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type actorContextKey struct{}

var ErrNoActor = errors.New("no authenticated user in context")

func WithActor(ctx context.Context, actor identity.UserID) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func actorFrom(ctx context.Context) (domain.UserID, error) {
	actor, ok := ctx.Value(actorContextKey{}).(identity.UserID)
	if !ok {
		return domain.UserID{}, ErrNoActor
	}

	return domain.UserID(actor), nil
}
