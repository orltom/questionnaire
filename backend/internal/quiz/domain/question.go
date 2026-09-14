package domain

import (
	"errors"
	"slices"
	"uuid"
)

type QuestionID uuid.UUID

type Question struct {
	id          QuestionID
	description string
	answers     []Answer
}

func NewQuestion(description string) (Question, error) {
	if len(description) == 0 {
		return Question{}, errors.New("description can not be empty")
	}

	return Question{
		id:          QuestionID(uuid.NewV7()),
		description: description,
		answers:     []Answer{},
	}, nil
}

type Answer struct {
	description string
	position    int
	count       int
}

func (q *Question) ID() QuestionID {
	return q.id
}

func (q *Question) Description() string {
	return q.description
}

func (q *Question) Answers() []Answer {
	return slices.Clone(q.answers)
}

func (q *Question) ChangeDescription(description string) error {
	if len(description) == 0 {
		return errors.New("description can not be empty")
	}

	q.description = description

	return nil
}

func (q *Question) ClearAnswers() {
	q.answers = []Answer{}
}

func (q *Question) AddAnswer(description string, position int, count int) error {
	if position < 0 {
		return errors.New("position can not be negative")
	}

	if len(description) == 0 {
		return errors.New("description can not be empty")
	}

	q.answers = append(q.answers, Answer{
		description: description,
		position:    position,
		count:       count,
	})

	return nil
}

func (a Answer) Description() string {
	return a.description
}

func (a Answer) Position() int {
	return a.position
}

func (a Answer) Count() int {
	return a.count
}
