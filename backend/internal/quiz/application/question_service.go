package application

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type QuestionRepository interface {
	Create(ctx context.Context, question domain.Question) error
	Find(ctx context.Context, id domain.QuestionID) (domain.Question, error)
	Update(ctx context.Context, question domain.Question) error
	Delete(ctx context.Context, id domain.QuestionID) error
}

type QuestionService struct {
	repository QuestionRepository
}

func NewQuestionService(repository QuestionRepository) *QuestionService {
	return &QuestionService{
		repository: repository,
	}
}

type QuestionRequest struct {
	Description string
	Answer      []AnswerRequest
}

type AnswerRequest struct {
	Description string
	Position    int
	Count       int
}

func (s *QuestionService) Create(ctx context.Context, req QuestionRequest) (domain.Question, error) {
	question, err := domain.NewQuestion(req.Description)
	if err != nil {
		return question, fmt.Errorf("%w: create question: %w", ErrInvalidArguments, err)
	}

	for _, a := range req.Answer {
		err = question.AddAnswer(a.Description, a.Position, a.Count)
		if err != nil {
			return question, fmt.Errorf("%w: add answer: %w", ErrInvalidArguments, err)
		}
	}

	err = s.repository.Create(ctx, question)
	if err != nil {
		return question, fmt.Errorf("%w: save question: %w", ErrPersistence, err)
	}

	return question, nil
}

func (s *QuestionService) Get(ctx context.Context, id domain.QuestionID) (domain.Question, error) {
	question, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return domain.Question{}, ErrEntityNotFound
		}

		return question, fmt.Errorf("%w: load question: %w", ErrPersistence, err)
	}

	return question, nil
}

func (s *QuestionService) Update(ctx context.Context, id domain.QuestionID, req QuestionRequest) error {
	question, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: load question by id '%s': %w", ErrPersistence, id, err)
	}

	err = question.ChangeDescription(req.Description)
	if err != nil {
		return fmt.Errorf("%w: update question description: %w", ErrInvalidArguments, err)
	}

	question.ClearAnswers()

	for _, a := range req.Answer {
		err := question.AddAnswer(a.Description, a.Position, a.Count)
		if err != nil {
			return fmt.Errorf("%w: add answer to question: %w", ErrInvalidArguments, err)
		}
	}

	err = s.repository.Update(ctx, question)
	if err != nil {
		return fmt.Errorf("%w: save question: %w", ErrPersistence, err)
	}

	return nil
}

func (s *QuestionService) Delete(ctx context.Context, id domain.QuestionID) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: delete question: %w", ErrPersistence, err)
	}

	return nil
}
