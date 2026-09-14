package domain

import (
	"errors"
	"slices"
	"uuid"
)

type QuizID uuid.UUID

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type Quiz struct {
	id          QuizID
	title       string
	description string
	visibility  Visibility
	questions   []QuizQuestion
}

func NewQuiz(title string, description string, visibility Visibility, questions ...QuizQuestion) (Quiz, error) {
	if len(title) == 0 {
		return Quiz{}, errors.New("quiz title is required")
	}

	return Quiz{
		id:          QuizID(uuid.NewV7()),
		title:       title,
		description: description,
		visibility:  visibility,
		questions:   questions,
	}, nil
}

type QuizQuestion struct {
	id       QuestionID
	position int
}

func (q *Quiz) ID() QuizID {
	return q.id
}

func (q *Quiz) Title() string {
	return q.title
}

func (q *Quiz) Description() string {
	return q.description
}

func (q *Quiz) Visibility() Visibility {
	return q.visibility
}

func (q *Quiz) Questions() []QuizQuestion {
	return slices.Clone(q.questions)
}

func (q *Quiz) Rename(title string) error {
	if len(title) == 0 {
		return errors.New("quiz title is required")
	}

	q.title = title

	return nil
}

func (q *Quiz) ChangeDescription(description string) {
	q.description = description
}

func (q *Quiz) ChangeVisibility(visibility Visibility) {
	q.visibility = visibility
}

func (q *Quiz) AddQuestion(id QuestionID, position int) error {
	if position < 0 {
		return errors.New("invalid position")
	}

	q.questions = append(q.questions, QuizQuestion{
		id:       id,
		position: position,
	})

	return nil
}

func (q QuizQuestion) ID() QuestionID {
	return q.id
}

func (q QuizQuestion) Position() int {
	return q.position
}
