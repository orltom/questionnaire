package application

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type QuizRepository interface {
	Create(ctx context.Context, quiz domain.Quiz) error
	Find(ctx context.Context, id domain.QuizID) (domain.Quiz, error)
	Update(ctx context.Context, quiz domain.Quiz) error
	Delete(ctx context.Context, id domain.QuizID) error
}

type QuizService struct {
	repository    QuizRepository
	lookupService QuestionLookupService
}

func NewQuizService(repository QuizRepository, lookupService QuestionLookupService) *QuizService {
	return &QuizService{
		repository:    repository,
		lookupService: lookupService,
	}
}

func (s *QuizService) Create(ctx context.Context, title string, description string) (domain.Quiz, error) {
	quiz, err := domain.NewQuiz(title, description, domain.Private)
	if err != nil {
		return quiz, fmt.Errorf("%w: create quiz: %w", ErrInvalidArguments, err)
	}

	err = s.repository.Create(ctx, quiz)
	if err != nil {
		return quiz, fmt.Errorf("%w: save quiz: %w", ErrPersistence, err)
	}

	return quiz, nil
}

func (s *QuizService) Get(ctx context.Context, id domain.QuizID) (domain.Quiz, error) {
	quiz, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return domain.Quiz{}, ErrEntityNotFound
		}

		return quiz, fmt.Errorf("%w: load quiz: %w", ErrPersistence, err)
	}

	return quiz, nil
}

func (s *QuizService) Update(
	ctx context.Context,
	id domain.QuizID,
	title string,
	description string,
	visibility domain.Visibility,
) error {
	quiz, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: find quiz: %w", ErrPersistence, err)
	}

	err = quiz.Rename(title)
	if err != nil {
		return fmt.Errorf("%w: update quiz title: %w", ErrInvalidArguments, err)
	}

	quiz.ChangeDescription(description)
	quiz.ChangeVisibility(visibility)

	err = s.repository.Update(ctx, quiz)
	if err != nil {
		return fmt.Errorf("%w: save quiz: %w", ErrPersistence, err)
	}

	return nil
}

func (s *QuizService) Delete(ctx context.Context, id domain.QuizID) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: delete quiz: %w", ErrPersistence, err)
	}

	return nil
}

func (s *QuizService) AddQuestion(ctx context.Context, id domain.QuizID, qID domain.QuestionID, pos int) error {
	ok, err := s.lookupService.Exists(ctx, qID)
	if err != nil {
		return fmt.Errorf("%w: lookup question: %w", ErrPersistence, err)
	}

	if !ok {
		return fmt.Errorf("%w: question with ID '%s' does not exist", ErrInvalidArguments, qID)
	}

	quiz, err := s.repository.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: could not find quiz with ID '%s': %w", ErrPersistence, id, err)
	}

	err = quiz.AddQuestion(qID, pos)
	if err != nil {
		return fmt.Errorf(
			"%w: could not add question with ID '%s' to quiz with ID '%s': %w",
			ErrInvalidArguments, qID, id, err,
		)
	}

	err = s.repository.Update(ctx, quiz)
	if err != nil {
		return fmt.Errorf("%w: save quiz: %w", ErrPersistence, err)
	}

	return nil
}
