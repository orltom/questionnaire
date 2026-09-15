package domain

import (
	"errors"
	"slices"
	"uuid"
)

type QuizID uuid.UUID

type Quiz struct {
	id          QuizID
	title       string
	description string
	owner       UserID
	visibility  Visibility
	questions   []QuizQuestion
}

func NewQuiz(title string, description string, owner UserID, visibility Visibility, questions ...QuizQuestion) (Quiz, error) {
	if len(title) == 0 {
		return Quiz{}, errors.New("quiz title is required")
	}

	if !visibility.Valid() {
		return Quiz{}, errors.New("unknown visibility")
	}

	return Quiz{
		id:          QuizID(uuid.NewV7()),
		title:       title,
		description: description,
		owner:       owner,
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

func (q *Quiz) Owner() UserID {
	return q.owner
}

func (q *Quiz) Visibility() Visibility {
	return q.visibility
}

func (q *Quiz) Questions() []QuizQuestion {
	return slices.Clone(q.questions)
}

func (q *Quiz) CanView(actor UserID) bool {
	return q.visibility == Public || q.owner == actor
}

func (q *Quiz) CanEdit(actor UserID) bool {
	return q.owner == actor
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

func (q *Quiz) ChangeVisibility(visibility Visibility) error {
	if !visibility.Valid() {
		return errors.New("unknown visibility")
	}

	q.visibility = visibility

	return nil
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
