package application

import (
	"context"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type QuestionLookupService interface {
	Exists(ctx context.Context, id domain.QuestionID) (bool, error)
}
