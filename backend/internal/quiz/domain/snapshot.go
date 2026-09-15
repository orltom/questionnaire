package domain

import (
	"errors"
	"slices"
)

type SnapshotAnswer struct {
	id          AnswerID
	description string
	position    int
	correct     bool
}

func NewSnapshotAnswer(id AnswerID, description string, position int, correct bool) SnapshotAnswer {
	return SnapshotAnswer{
		id:          id,
		description: description,
		position:    position,
		correct:     correct,
	}
}

func (a SnapshotAnswer) ID() AnswerID {
	return a.id
}

func (a SnapshotAnswer) Description() string {
	return a.description
}

func (a SnapshotAnswer) Position() int {
	return a.position
}

func (a SnapshotAnswer) Correct() bool {
	return a.correct
}

type SnapshotQuestion struct {
	id          QuestionID
	description string
	position    int
	answers     []SnapshotAnswer
}

func NewSnapshotQuestion(id QuestionID, description string, position int, answers ...SnapshotAnswer) SnapshotQuestion {
	return SnapshotQuestion{
		id:          id,
		description: description,
		position:    position,
		answers:     answers,
	}
}

func (q SnapshotQuestion) ID() QuestionID {
	return q.id
}

func (q SnapshotQuestion) Description() string {
	return q.description
}

func (q SnapshotQuestion) Position() int {
	return q.position
}

func (q SnapshotQuestion) Answers() []SnapshotAnswer {
	return slices.Clone(q.answers)
}

type QuizSnapshot struct {
	quizID    QuizID
	title     string
	questions []SnapshotQuestion
}

func NewQuizSnapshot(quizID QuizID, title string, questions ...SnapshotQuestion) (QuizSnapshot, error) {
	if len(title) == 0 {
		return QuizSnapshot{}, errors.New("quiz title is required")
	}

	if len(questions) == 0 {
		return QuizSnapshot{}, errors.New("quiz without questions can not be challenged")
	}

	return QuizSnapshot{
		quizID:    quizID,
		title:     title,
		questions: questions,
	}, nil
}

func (s QuizSnapshot) QuizID() QuizID {
	return s.quizID
}

func (s QuizSnapshot) Title() string {
	return s.title
}

func (s QuizSnapshot) Questions() []SnapshotQuestion {
	return slices.Clone(s.questions)
}

func (s QuizSnapshot) Contains(questionID QuestionID, answerID AnswerID) bool {
	for _, q := range s.questions {
		if q.id != questionID {
			continue
		}

		for _, a := range q.answers {
			if a.id == answerID {
				return true
			}
		}
	}

	return false
}

func (s QuizSnapshot) Correct(answerID AnswerID) bool {
	for _, q := range s.questions {
		for _, a := range q.answers {
			if a.id == answerID {
				return a.correct
			}
		}
	}

	return false
}
