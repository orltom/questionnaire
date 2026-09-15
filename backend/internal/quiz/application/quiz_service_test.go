package application

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestQuizService_Create(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	defaultQuiz, _ := domain.NewQuiz("Animals", "Quiz about animals", actor, domain.Private)

	type mocks struct {
		repository func(r *MockQuizRepository)
	}
	type args struct {
		title       string
		description string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    domain.Quiz
		wantErr error
	}{
		{
			name: "When creating a valid quiz, then save it",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Cond(func(q domain.Quiz) bool {
						return q.Title() == "Animals" && q.Description() == "Quiz about animals" && q.Visibility() == domain.Private
					})).Return(nil)
				},
			},

			args: args{
				title:       "Animals",
				description: "Quiz about animals",
			},
			want:    defaultQuiz,
			wantErr: nil,
		},
		{
			name: "When the title is empty, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {},
			},
			args: args{
				title:       "",
				description: "",
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				title:       "Animals",
				description: "Quiz about animals",
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockQuizRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuizService{repository: repo, lookupService: nil}
			got, err := s.Create(context.Background(), actor, tt.args.title, tt.args.description, domain.Private)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(domain.Quiz{}, domain.QuizQuestion{}), cmpopts.IgnoreFields(domain.Quiz{}, "id")); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestQuizService_Get(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	defaultQuiz, _ := domain.NewQuiz("Animals", "Quiz about animals", actor, domain.Private)

	type mocks struct {
		repository func(r *MockQuizRepository)
	}
	type args struct {
		id domain.QuizID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    domain.Quiz
		wantErr error
	}{
		{
			name: "When the quiz exists, then return it",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(defaultQuiz, nil)
				},
			},

			args: args{
				id: defaultQuiz.ID(),
			},
			want:    defaultQuiz,
			wantErr: nil,
		},
		{
			name: "When the quiz ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(domain.Quiz{}, ErrEntityNotFound)
				},
			},
			args: args{
				id: defaultQuiz.ID(),
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(domain.Quiz{}, errors.New("database is down"))
				},
			},
			args: args{
				id: defaultQuiz.ID(),
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockQuizRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuizService{repository: repo, lookupService: nil}
			got, err := s.Get(context.Background(), actor, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(domain.Quiz{}, domain.QuizQuestion{})); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestQuizService_Update(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	defaultQuiz, _ := domain.NewQuiz("Animals", "Quiz about animals", actor, domain.Private)

	type mocks struct {
		repository func(r *MockQuizRepository)
	}
	type args struct {
		id          domain.QuizID
		title       string
		description string
		visibility  domain.Visibility
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When updating an existing quiz, then save the changes",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(defaultQuiz, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Cond(func(q domain.Quiz) bool {
						return q.Title() == "Sport" &&
							q.Description() == "Chess is a sport" &&
							q.Visibility() == domain.Public
					})).Return(nil)
				},
			},

			args: args{
				id:          defaultQuiz.ID(),
				title:       "Sport",
				description: "Chess is a sport",
				visibility:  domain.Public,
			},
			wantErr: nil,
		},
		{
			name: "When the quiz ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(domain.Quiz{}, ErrEntityNotFound)
				},
			},
			args: args{
				id: defaultQuiz.ID(),
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(domain.Quiz{}, errors.New("database is down"))
				},
			},
			args: args{
				id: defaultQuiz.ID(),
			},
			wantErr: ErrPersistence,
		},
		{
			name: "When the title is empty, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(defaultQuiz, nil)
				},
			},
			args: args{
				id:          defaultQuiz.ID(),
				title:       "",
				description: "Chess is a sport",
				visibility:  domain.Public,
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When saving the quiz fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuiz.ID()).Return(defaultQuiz, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				id:          defaultQuiz.ID(),
				title:       "Sport",
				description: "Chess is a sport",
				visibility:  domain.Public,
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockQuizRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuizService{repository: repo, lookupService: nil}
			err := s.Update(context.Background(), actor, tt.args.id, tt.args.title, tt.args.description, tt.args.visibility)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizService_Delete(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	id := domain.QuizID(uuid.Nil())
	defaultQuiz, _ := domain.NewQuiz("Animals", "Quiz about animals", actor, domain.Private)

	type mocks struct {
		repository func(r *MockQuizRepository)
	}
	type args struct {
		id domain.QuizID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When the quiz exists, then delete it",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), id).Return(defaultQuiz, nil)
					r.EXPECT().Delete(gomock.Any(), id).Return(nil)
				},
			},

			args: args{
				id: id,
			},
			wantErr: nil,
		},
		{
			name: "When the quiz ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), id).Return(domain.Quiz{}, ErrEntityNotFound)
				},
			},
			args: args{
				id: id,
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuizRepository) {
					r.EXPECT().Find(gomock.Any(), id).Return(defaultQuiz, nil)
					r.EXPECT().Delete(gomock.Any(), id).Return(errors.New("database is down"))
				},
			},
			args: args{
				id: id,
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockQuizRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuizService{repository: repo, lookupService: nil}
			err := s.Delete(context.Background(), actor, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
