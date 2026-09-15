package domain

import (
	"slices"
	"uuid"
)

type ParticipationID uuid.UUID

type GivenAnswer struct {
	questionID QuestionID
	answerID   AnswerID
}

func (g GivenAnswer) QuestionID() QuestionID {
	return g.questionID
}

func (g GivenAnswer) AnswerID() AnswerID {
	return g.answerID
}

type Participation struct {
	id          ParticipationID
	challengeID ChallengeID
	userID      UserID
	answers     []GivenAnswer
}

func NewParticipation(challengeID ChallengeID, userID UserID) Participation {
	return Participation{
		id:          ParticipationID(uuid.NewV7()),
		challengeID: challengeID,
		userID:      userID,
		answers:     []GivenAnswer{},
	}
}

func (p *Participation) ID() ParticipationID {
	return p.id
}

func (p *Participation) ChallengeID() ChallengeID {
	return p.challengeID
}

func (p *Participation) UserID() UserID {
	return p.userID
}

func (p *Participation) Answers() []GivenAnswer {
	return slices.Clone(p.answers)
}

func (p *Participation) Answer(questionID QuestionID, answerID AnswerID) {
	given := GivenAnswer{
		questionID: questionID,
		answerID:   answerID,
	}

	for i := range p.answers {
		if p.answers[i].questionID == questionID {
			p.answers[i] = given

			return
		}
	}

	p.answers = append(p.answers, given)
}

func (p *Participation) Points(snapshot QuizSnapshot) int {
	points := 0

	for _, a := range p.answers {
		if snapshot.Correct(a.answerID) {
			points++
		}
	}

	return points
}
