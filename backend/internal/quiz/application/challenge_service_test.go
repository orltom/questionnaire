package application

import (
	"context"
	"errors"
	"testing"
	"time"
	"uuid"

	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestChallengeService_Create(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	stranger := domain.UserID(uuid.NewV7())
	question, _ := domain.NewQuestion("Largest animal?", actor, domain.Public)
	_ = question.AddAnswer("Blue whale", 0, true)

	quiz, _ := domain.NewQuiz("Animals", "Questions about animals", actor, domain.Public)
	_ = quiz.AddQuestion(question.ID(), 0)

	foreignQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", stranger, domain.Private)

	type mocks struct {
		repository func(r *MockChallengeRepository)
		quizzes    func(r *MockQuizFinder)
		questions  func(r *MockQuestionFinder)
	}
	type args struct {
		actor   domain.UserID
		request ChallengeRequest
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When the quiz can be viewed, then snapshot it and save the challenge",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				},
				quizzes: func(r *MockQuizFinder) {
					r.EXPECT().Find(gomock.Any(), quiz.ID()).Return(quiz, nil)
				},
				questions: func(r *MockQuestionFinder) {
					r.EXPECT().Find(gomock.Any(), question.ID()).Return(question, nil)
				},
			},
			args: args{
				actor: actor,
				request: ChallengeRequest{
					QuizID:   quiz.ID(),
					StartsAt: time.Now(),
					EndsAt:   nil,
					Access:   domain.AccessPublic,
					Invitees: nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "When the quiz does not exist, then return an entity not found error",
			mocks: mocks{
				repository: func(_ *MockChallengeRepository) {},
				quizzes: func(r *MockQuizFinder) {
					r.EXPECT().Find(gomock.Any(), quiz.ID()).Return(domain.Quiz{}, ErrEntityNotFound)
				},
				questions: func(_ *MockQuestionFinder) {},
			},
			args: args{
				actor: actor,
				request: ChallengeRequest{
					QuizID:   quiz.ID(),
					StartsAt: time.Now(),
					EndsAt:   nil,
					Access:   domain.AccessPublic,
					Invitees: nil,
				},
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the quiz is private and owned by someone else, then return a forbidden error",
			mocks: mocks{
				repository: func(_ *MockChallengeRepository) {},
				quizzes: func(r *MockQuizFinder) {
					r.EXPECT().Find(gomock.Any(), foreignQuiz.ID()).Return(foreignQuiz, nil)
				},
				questions: func(_ *MockQuestionFinder) {},
			},
			args: args{
				actor: actor,
				request: ChallengeRequest{
					QuizID:   foreignQuiz.ID(),
					StartsAt: time.Now(),
					EndsAt:   nil,
					Access:   domain.AccessPublic,
					Invitees: nil,
				},
			},
			wantErr: ErrForbidden,
		},
		{
			name: "When a referenced question is gone, then return an invalid arguments error",
			mocks: mocks{
				repository: func(_ *MockChallengeRepository) {},
				quizzes: func(r *MockQuizFinder) {
					r.EXPECT().Find(gomock.Any(), quiz.ID()).Return(quiz, nil)
				},
				questions: func(r *MockQuestionFinder) {
					r.EXPECT().Find(gomock.Any(), question.ID()).Return(domain.Question{}, ErrEntityNotFound)
				},
			},
			args: args{
				actor: actor,
				request: ChallengeRequest{
					QuizID:   quiz.ID(),
					StartsAt: time.Now(),
					EndsAt:   nil,
					Access:   domain.AccessPublic,
					Invitees: nil,
				},
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the access is unknown, then return an invalid arguments error",
			mocks: mocks{
				repository: func(_ *MockChallengeRepository) {},
				quizzes: func(r *MockQuizFinder) {
					r.EXPECT().Find(gomock.Any(), quiz.ID()).Return(quiz, nil)
				},
				questions: func(r *MockQuestionFinder) {
					r.EXPECT().Find(gomock.Any(), question.ID()).Return(question, nil)
				},
			},
			args: args{
				actor: actor,
				request: ChallengeRequest{
					QuizID:   quiz.ID(),
					StartsAt: time.Now(),
					EndsAt:   nil,
					Access:   domain.Access("everyone"),
					Invitees: nil,
				},
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
				quizzes: func(r *MockQuizFinder) {
					r.EXPECT().Find(gomock.Any(), quiz.ID()).Return(quiz, nil)
				},
				questions: func(r *MockQuestionFinder) {
					r.EXPECT().Find(gomock.Any(), question.ID()).Return(question, nil)
				},
			},
			args: args{
				actor: actor,
				request: ChallengeRequest{
					QuizID:   quiz.ID(),
					StartsAt: time.Now(),
					EndsAt:   nil,
					Access:   domain.AccessPublic,
					Invitees: nil,
				},
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			challenges := NewMockChallengeRepository(ctrl)
			quizzes := NewMockQuizFinder(ctrl)
			questions := NewMockQuestionFinder(ctrl)
			tt.mocks.repository(challenges)
			tt.mocks.quizzes(quizzes)
			tt.mocks.questions(questions)

			s := NewChallengeService(challenges, quizzes, questions)
			_, err := s.Create(context.Background(), tt.args.actor, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChallengeService_Get(t *testing.T) {
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

	type mocks struct {
		repository func(r *MockChallengeRepository)
	}
	type args struct {
		actor domain.UserID
		id    domain.ChallengeID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When the challenge is public, then anyone can load it",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(publicChallenge, nil)
				},
			},
			args: args{
				actor: stranger,
				id:    publicChallenge.ID(),
			},
			wantErr: nil,
		},
		{
			name: "When the challenge is invite only, then a stranger is not allowed",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
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
			name: "When the challenge is invite only, then the owner is allowed",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), invitedChallenge.ID()).Return(invitedChallenge, nil)
				},
			},
			args: args{
				actor: actor,
				id:    invitedChallenge.ID(),
			},
			wantErr: nil,
		},
		{
			name: "When the challenge ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
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
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), publicChallenge.ID()).Return(domain.Challenge{}, errors.New("database is down"))
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
			challenges := NewMockChallengeRepository(ctrl)
			tt.mocks.repository(challenges)

			s := NewChallengeService(challenges, NewMockQuizFinder(ctrl), NewMockQuestionFinder(ctrl))
			_, err := s.Get(context.Background(), tt.args.actor, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChallengeService_Close(t *testing.T) {
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

	challenge, _ := domain.NewChallenge(actor, snapshot, time.Now().Add(-time.Hour), nil, domain.AccessPublic)

	type mocks struct {
		repository func(r *MockChallengeRepository)
	}
	type args struct {
		actor domain.UserID
		id    domain.ChallengeID
		now   time.Time
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When the owner closes the challenge, then save the changes",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), challenge.ID()).Return(challenge, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			args: args{
				actor: actor,
				id:    challenge.ID(),
				now:   time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "When a stranger closes the challenge, then return a forbidden error",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), challenge.ID()).Return(challenge, nil)
				},
			},
			args: args{
				actor: stranger,
				id:    challenge.ID(),
				now:   time.Now(),
			},
			wantErr: ErrForbidden,
		},
		{
			name: "When the challenge ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), challenge.ID()).Return(domain.Challenge{}, ErrEntityNotFound)
				},
			},
			args: args{
				actor: actor,
				id:    challenge.ID(),
				now:   time.Now(),
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When saving the challenge fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockChallengeRepository) {
					r.EXPECT().Find(gomock.Any(), challenge.ID()).Return(challenge, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				actor: actor,
				id:    challenge.ID(),
				now:   time.Now(),
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			challenges := NewMockChallengeRepository(ctrl)
			tt.mocks.repository(challenges)

			s := NewChallengeService(challenges, NewMockQuizFinder(ctrl), NewMockQuestionFinder(ctrl))
			err := s.Close(context.Background(), tt.args.actor, tt.args.id, tt.args.now)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
