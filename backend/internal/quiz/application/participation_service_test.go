package application

import (
	"context"
	"errors"
	"testing"
	"time"
	"uuid"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestParticipationService_Submit(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	stranger := domain.UserID(uuid.NewV7())
	now := time.Now()

	question, _ := domain.NewQuestion("Largest animal?", actor, domain.Public)
	_ = question.AddAnswer("Blue whale", 0, true)
	answer := question.Answers()[0]

	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(
			question.ID(),
			question.Description(),
			0,
			domain.NewSnapshotAnswer(answer.ID(), answer.Description(), answer.Position(), answer.Correct()),
		),
	)

	publicChallenge, _ := domain.NewChallenge(actor, snapshot, now.Add(-time.Hour), nil, domain.AccessPublic)
	invitedChallenge, _ := domain.NewChallenge(actor, snapshot, now.Add(-time.Hour), nil, domain.AccessInvited)
	participation := domain.NewParticipation(publicChallenge.ID(), actor)

	type mocks struct {
		repository func(r *MockParticipationRepository)
		challenges func(r *MockChallengeRepository)
	}
	type args struct {
		actor      domain.UserID
		id         domain.ChallengeID
		questionID domain.QuestionID
		answerID   domain.AnswerID
		now        time.Time
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When the user has not answered yet, then create a participation",
			mocks: mocks{
				repository: func(r *MockParticipationRepository) {
					r.EXPECT().FindByUser(gomock.Any(), publicChallenge.ID(), actor).Return(domain.Participation{}, ErrEntityNotFound)
					r.EXPECT().Create(gomock.Any(), gomock.Cond(func(p domain.Participation) bool {
						return len(p.Answers()) == 1 && p.Answers()[0].AnswerID() == answer.ID()
					})).Return(nil)
				},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(publicChallenge, nil)
				},
			},
			args: args{
				actor:      actor,
				id:         publicChallenge.ID(),
				questionID: question.ID(),
				answerID:   answer.ID(),
				now:        now,
			},
			wantErr: nil,
		},
		{
			name: "When the user has already participated, then update the participation",
			mocks: mocks{
				repository: func(r *MockParticipationRepository) {
					r.EXPECT().FindByUser(gomock.Any(), publicChallenge.ID(), actor).Return(participation, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(publicChallenge, nil)
				},
			},
			args: args{
				actor:      actor,
				id:         publicChallenge.ID(),
				questionID: question.ID(),
				answerID:   answer.ID(),
				now:        now,
			},
			wantErr: nil,
		},
		{
			name: "When the answer is not part of the quiz, then return an invalid arguments error",
			mocks: mocks{
				repository: func(_ *MockParticipationRepository) {},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(publicChallenge, nil)
				},
			},
			args: args{
				actor:      actor,
				id:         publicChallenge.ID(),
				questionID: domain.QuestionID{},
				answerID:   domain.AnswerID{},
				now:        now,
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the challenge is invite only, then a stranger can not participate",
			mocks: mocks{
				repository: func(_ *MockParticipationRepository) {},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), invitedChallenge.ID()).Return(invitedChallenge, nil)
				},
			},
			args: args{
				actor:      stranger,
				id:         invitedChallenge.ID(),
				questionID: question.ID(),
				answerID:   answer.ID(),
				now:        now,
			},
			wantErr: ErrForbidden,
		},
		{
			name: "When the challenge ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(_ *MockParticipationRepository) {},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(domain.Challenge{}, ErrEntityNotFound)
				},
			},
			args: args{
				actor:      actor,
				id:         publicChallenge.ID(),
				questionID: question.ID(),
				answerID:   answer.ID(),
				now:        now,
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When saving the participation fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockParticipationRepository) {
					r.EXPECT().FindByUser(gomock.Any(), publicChallenge.ID(), actor).Return(domain.Participation{}, ErrEntityNotFound)
					r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(publicChallenge, nil)
				},
			},
			args: args{
				actor:      actor,
				id:         publicChallenge.ID(),
				questionID: question.ID(),
				answerID:   answer.ID(),
				now:        now,
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockParticipationRepository(ctrl)
			challenges := NewMockChallengeRepository(ctrl)
			tt.mocks.repository(repo)
			tt.mocks.challenges(challenges)

			s := NewParticipationService(repo, challenges)
			err := s.Submit(context.Background(), tt.args.actor, tt.args.id, tt.args.questionID, tt.args.answerID, tt.args.now)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Submit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParticipationService_Results(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	stranger := domain.UserID(uuid.NewV7())
	question, _ := domain.NewQuestion("Largest animal?", actor, domain.Public)
	_ = question.AddAnswer("Blue whale", 0, true)
	answer := question.Answers()[0]

	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(
			question.ID(),
			question.Description(),
			0,
			domain.NewSnapshotAnswer(answer.ID(), answer.Description(), answer.Position(), answer.Correct()),
		),
	)

	publicChallenge, _ := domain.NewChallenge(actor, snapshot, time.Now().Add(-time.Hour), nil, domain.AccessPublic)
	invitedChallenge, _ := domain.NewChallenge(actor, snapshot, time.Now().Add(-time.Hour), nil, domain.AccessInvited)

	participation := domain.NewParticipation(publicChallenge.ID(), actor)
	participation.Answer(question.ID(), answer.ID())

	type mocks struct {
		repository func(r *MockParticipationRepository)
		challenges func(r *MockChallengeRepository)
	}
	type args struct {
		actor domain.UserID
		id    domain.ChallengeID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    []Result
		wantErr error
	}{
		{
			name: "When the challenge has participations, then return the points per user",
			mocks: mocks{
				repository: func(r *MockParticipationRepository) {
					r.EXPECT().FindByChallenge(gomock.Any(), publicChallenge.ID()).Return([]domain.Participation{participation}, nil)
				},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(publicChallenge, nil)
				},
			},
			args: args{
				actor: actor,
				id:    publicChallenge.ID(),
			},
			want:    []Result{{User: actor, Points: 1}},
			wantErr: nil,
		},
		{
			name: "When the challenge is invite only, then a stranger can not see the results",
			mocks: mocks{
				repository: func(_ *MockParticipationRepository) {},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), invitedChallenge.ID()).Return(invitedChallenge, nil)
				},
			},
			args: args{
				actor: stranger,
				id:    invitedChallenge.ID(),
			},
			wantErr: ErrForbidden,
		},
		{
			name: "When the challenge ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(_ *MockParticipationRepository) {},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(domain.Challenge{}, ErrEntityNotFound)
				},
			},
			args: args{
				actor: actor,
				id:    publicChallenge.ID(),
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockParticipationRepository) {
					r.EXPECT().FindByChallenge(gomock.Any(), publicChallenge.ID()).Return(nil, errors.New("database is down"))
				},
				challenges: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(publicChallenge, nil)
				},
			},
			args: args{
				actor: actor,
				id:    publicChallenge.ID(),
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockParticipationRepository(ctrl)
			challenges := NewMockChallengeRepository(ctrl)
			tt.mocks.repository(repo)
			tt.mocks.challenges(challenges)

			s := NewParticipationService(repo, challenges)
			got, err := s.Results(context.Background(), tt.args.actor, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Results() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Results() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
