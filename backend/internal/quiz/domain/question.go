package domain

import (
	"errors"
	"slices"
	"uuid"
)

type QuestionID uuid.UUID

type AnswerID uuid.UUID

type Label string

type Question struct {
	id          QuestionID
	description string
	owner       UserID
	visibility  Visibility
	labels      []Label
	answers     []Answer
}

func NewQuestion(description string, owner UserID, visibility Visibility) (Question, error) {
	if len(description) == 0 {
		return Question{}, errors.New("description can not be empty")
	}

	if !visibility.Valid() {
		return Question{}, errors.New("unknown visibility")
	}

	return Question{
		id:          QuestionID(uuid.NewV7()),
		description: description,
		owner:       owner,
		visibility:  visibility,
		labels:      []Label{},
		answers:     []Answer{},
	}, nil
}

type Answer struct {
	id          AnswerID
	description string
	position    int
	correct     bool
}

func (q *Question) ID() QuestionID {
	return q.id
}

func (q *Question) Description() string {
	return q.description
}

func (q *Question) Owner() UserID {
	return q.owner
}

func (q *Question) Visibility() Visibility {
	return q.visibility
}

func (q *Question) Labels() []Label {
	return slices.Clone(q.labels)
}

func (q *Question) Answers() []Answer {
	return slices.Clone(q.answers)
}

func (q *Question) CanView(actor UserID) bool {
	return q.visibility == Public || q.owner == actor
}

func (q *Question) CanEdit(actor UserID) bool {
	return q.owner == actor
}

func (q *Question) ChangeDescription(description string) error {
	if len(description) == 0 {
		return errors.New("description can not be empty")
	}

	q.description = description

	return nil
}

func (q *Question) ChangeVisibility(visibility Visibility) error {
	if !visibility.Valid() {
		return errors.New("unknown visibility")
	}

	q.visibility = visibility

	return nil
}

func (q *Question) ClearLabels() {
	q.labels = []Label{}
}

func (q *Question) AddLabel(label Label) error {
	if len(label) == 0 {
		return errors.New("label can not be empty")
	}

	if slices.Contains(q.labels, label) {
		return nil
	}

	q.labels = append(q.labels, label)

	return nil
}

func (q *Question) ClearAnswers() {
	q.answers = []Answer{}
}

func (q *Question) AddAnswer(description string, position int, correct bool) error {
	if position < 0 {
		return errors.New("position can not be negative")
	}

	if len(description) == 0 {
		return errors.New("description can not be empty")
	}

	q.answers = append(q.answers, Answer{
		id:          AnswerID(uuid.NewV7()),
		description: description,
		position:    position,
		correct:     correct,
	})

	return nil
}

func (a Answer) ID() AnswerID {
	return a.id
}

func (a Answer) Description() string {
	return a.description
}

func (a Answer) Position() int {
	return a.position
}

func (a Answer) Correct() bool {
	return a.correct
}
