package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type ChallengeRepository interface {
	Create(ctx context.Context, challenge domain.Challenge) error
	Find(ctx context.Context, id domain.ChallengeID) (domain.Challenge, error)
	Update(ctx context.Context, challenge domain.Challenge) error
	Delete(ctx context.Context, id domain.ChallengeID) error
}

type QuizFinder interface {
	Find(ctx context.Context, id domain.QuizID) (domain.Quiz, error)
}

type QuestionFinder interface {
	Find(ctx context.Context, id domain.QuestionID) (domain.Question, error)
}

type ChallengeService struct {
	repository      ChallengeRepository
	quizFinder      QuizFinder
	questionsFinder QuestionFinder
}

func NewChallengeService(repository ChallengeRepository, quizzes QuizFinder, questions QuestionFinder) *ChallengeService {
	return &ChallengeService{
		repository:      repository,
		quizFinder:      quizzes,
		questionsFinder: questions,
	}
}

type ChallengeRequest struct {
	QuizID   domain.QuizID
	StartsAt time.Time
	EndsAt   *time.Time
	Access   domain.Access
	Invitees []domain.UserID
}

func (s *ChallengeService) Create(ctx context.Context, actor domain.UserID, req ChallengeRequest) (domain.Challenge, error) {
	quiz, err := s.quizFinder.Find(ctx, req.QuizID)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return domain.Challenge{}, ErrEntityNotFound
		}

		return domain.Challenge{}, fmt.Errorf("%w: load quiz: %w", ErrPersistence, err)
	}

	if !quiz.CanView(actor) {
		return domain.Challenge{}, ErrForbidden
	}

	snapshot, err := s.snapshot(ctx, quiz)
	if err != nil {
		return domain.Challenge{}, err
	}

	challenge, err := domain.NewChallenge(actor, snapshot, req.StartsAt, req.EndsAt, req.Access, req.Invitees...)
	if err != nil {
		return domain.Challenge{}, fmt.Errorf("%w: create challenge: %w", ErrInvalidArguments, err)
	}

	err = s.repository.Create(ctx, challenge)
	if err != nil {
		return domain.Challenge{}, fmt.Errorf("%w: save challenge: %w", ErrPersistence, err)
	}

	return challenge, nil
}

func (s *ChallengeService) Get(ctx context.Context, actor domain.UserID, id domain.ChallengeID) (domain.Challenge, error) {
	challenge, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return domain.Challenge{}, ErrEntityNotFound
		}

		return domain.Challenge{}, fmt.Errorf("%w: load challenge: %w", ErrPersistence, err)
	}

	if !challenge.CanView(actor) {
		return domain.Challenge{}, ErrForbidden
	}

	return challenge, nil
}

func (s *ChallengeService) Close(ctx context.Context, actor domain.UserID, id domain.ChallengeID, now time.Time) error {
	challenge, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: find challenge: %w", ErrPersistence, err)
	}

	if !challenge.CanEdit(actor) {
		return ErrForbidden
	}

	err = challenge.Close(now)
	if err != nil {
		return fmt.Errorf("%w: close challenge: %w", ErrInvalidArguments, err)
	}

	err = s.repository.Update(ctx, challenge)
	if err != nil {
		return fmt.Errorf("%w: save challenge: %w", ErrPersistence, err)
	}

	return nil
}

func (s *ChallengeService) Delete(ctx context.Context, actor domain.UserID, id domain.ChallengeID) error {
	challenge, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: load challenge: %w", ErrPersistence, err)
	}

	if !challenge.CanEdit(actor) {
		return ErrForbidden
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: delete challenge: %w", ErrPersistence, err)
	}

	return nil
}

func (s *ChallengeService) snapshot(ctx context.Context, quiz domain.Quiz) (domain.QuizSnapshot, error) {
	questions := make([]domain.SnapshotQuestion, 0, len(quiz.Questions()))
	for _, ref := range quiz.Questions() {
		question, err := s.questionsFinder.Find(ctx, ref.ID())
		if err != nil {
			if errors.Is(err, ErrEntityNotFound) {
				return domain.QuizSnapshot{}, fmt.Errorf("%w: question '%s' is gone", ErrInvalidArguments, ref.ID())
			}

			return domain.QuizSnapshot{}, fmt.Errorf("%w: load question: %w", ErrPersistence, err)
		}

		answers := make([]domain.SnapshotAnswer, len(question.Answers()))
		for i, a := range question.Answers() {
			answers[i] = domain.NewSnapshotAnswer(a.ID(), a.Description(), a.Position(), a.Correct())
		}

		questions = append(
			questions,
			domain.NewSnapshotQuestion(question.ID(), question.Description(), ref.Position(), answers...),
		)
	}

	snapshot, err := domain.NewQuizSnapshot(quiz.ID(), quiz.Title(), questions...)
	if err != nil {
		return domain.QuizSnapshot{}, fmt.Errorf("%w: snapshot quiz: %w", ErrInvalidArguments, err)
	}

	return snapshot, nil
}
