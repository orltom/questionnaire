package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type ParticipationRepository interface {
	Create(ctx context.Context, participation domain.Participation) error
	Update(ctx context.Context, participation domain.Participation) error
	FindByUser(ctx context.Context, id domain.ChallengeID, user domain.UserID) (domain.Participation, error)
	FindByChallenge(ctx context.Context, id domain.ChallengeID) ([]domain.Participation, error)
}

type ParticipationService struct {
	repository ParticipationRepository
	challenges ChallengeRepository
}

func NewParticipationService(repository ParticipationRepository, challenges ChallengeRepository) *ParticipationService {
	return &ParticipationService{
		repository: repository,
		challenges: challenges,
	}
}

type Result struct {
	User   domain.UserID
	Points int
}

func (s *ParticipationService) Submit(ctx context.Context, actor domain.UserID, id domain.ChallengeID, questionID domain.QuestionID, answerID domain.AnswerID, now time.Time) error {
	challenge, err := s.challenges.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: load challenge: %w", ErrPersistence, err)
	}

	if !challenge.CanParticipate(actor, now) {
		return ErrForbidden
	}

	if !challenge.Quiz().Contains(questionID, answerID) {
		return fmt.Errorf("%w: answer does not belong to the challenged quiz", ErrInvalidArguments)
	}

	participation, err := s.repository.FindByUser(ctx, id, actor)
	if err != nil {
		if !errors.Is(err, ErrEntityNotFound) {
			return fmt.Errorf("%w: load participation: %w", ErrPersistence, err)
		}

		participation = domain.NewParticipation(id, actor)
		participation.Answer(questionID, answerID)

		err = s.repository.Create(ctx, participation)
		if err != nil {
			return fmt.Errorf("%w: save participation: %w", ErrPersistence, err)
		}

		return nil
	}

	participation.Answer(questionID, answerID)

	err = s.repository.Update(ctx, participation)
	if err != nil {
		return fmt.Errorf("%w: save participation: %w", ErrPersistence, err)
	}

	return nil
}

func (s *ParticipationService) Results(ctx context.Context, actor domain.UserID, id domain.ChallengeID) ([]Result, error) {
	challenge, err := s.challenges.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return nil, ErrEntityNotFound
		}

		return nil, fmt.Errorf("%w: load challenge: %w", ErrPersistence, err)
	}

	if !challenge.CanView(actor) {
		return nil, ErrForbidden
	}

	participations, err := s.repository.FindByChallenge(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: load participations: %w", ErrPersistence, err)
	}

	results := make([]Result, len(participations))
	for i, p := range participations {
		results[i] = Result{
			User:   p.UserID(),
			Points: p.Points(challenge.Quiz()),
		}
	}

	return results, nil
}
