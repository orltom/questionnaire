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

func TestQuestionService_Create(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	defaultQuestion, _ := domain.NewQuestion("What is the capital of Switzerland?", actor, domain.Private)
	_ = defaultQuestion.AddAnswer("Bern", 0, true)

	type mocks struct {
		repository func(r *MockQuestionRepository)
	}
	type args struct {
		req QuestionRequest
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    domain.Question
		wantErr error
	}{
		{
			name: "When creating a valid question, then save it",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Cond(func(q domain.Question) bool {
						answers := q.Answers()

						return q.Description() == "What is the capital of Switzerland?" &&
							len(answers) == 1 &&
							answers[0].Description() == "Bern"
					})).Return(nil)
				},
			},

			args: args{
				req: QuestionRequest{
					Description: "What is the capital of Switzerland?",
					Visibility:  domain.Private,
					Answer: []AnswerRequest{
						{Description: "Bern", Position: 0, Correct: true},
					},
				},
			},
			want:    defaultQuestion,
			wantErr: nil,
		},
		{
			name: "When the description is empty, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {},
			},
			args: args{
				req: QuestionRequest{
					Description: "",
					Visibility:  domain.Private,
					Answer:      nil,
				},
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When an answer position is negative, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {},
			},
			args: args{
				req: QuestionRequest{
					Description: "What is the capital of Switzerland?",
					Visibility:  domain.Private,
					Answer: []AnswerRequest{
						{Description: "Bern", Position: -1, Correct: true},
					},
				},
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				req: QuestionRequest{
					Description: "What is the capital of Switzerland?",
					Visibility:  domain.Private,
					Answer:      nil,
				},
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockQuestionRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuestionService{repository: repo}
			got, err := s.Create(context.Background(), actor, tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if diff := cmp.Diff(
				tt.want,
				got,
				cmp.AllowUnexported(domain.Question{}, domain.Answer{}),
				cmpopts.IgnoreFields(domain.Question{}, "id"),
				cmpopts.IgnoreFields(domain.Answer{}, "id"),
			); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestQuestionService_Get(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	defaultQuestion, _ := domain.NewQuestion("What is the capital of Switzerland?", actor, domain.Private)

	type mocks struct {
		repository func(r *MockQuestionRepository)
	}
	type args struct {
		id domain.QuestionID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    domain.Question
		wantErr error
	}{
		{
			name: "When the question exists, then return it",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(defaultQuestion, nil)
				},
			},

			args: args{
				id: defaultQuestion.ID(),
			},
			want:    defaultQuestion,
			wantErr: nil,
		},
		{
			name: "When the question ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(domain.Question{}, ErrEntityNotFound)
				},
			},
			args: args{
				id: defaultQuestion.ID(),
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(domain.Question{}, errors.New("database is down"))
				},
			},
			args: args{
				id: defaultQuestion.ID(),
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockQuestionRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuestionService{repository: repo}
			got, err := s.Get(context.Background(), actor, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(domain.Question{}, domain.Answer{})); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestQuestionService_Update(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	defaultQuestion, _ := domain.NewQuestion("What is the capital of Germany?", actor, domain.Private)
	_ = defaultQuestion.AddAnswer("Munich", 0, true)

	type mocks struct {
		repository func(r *MockQuestionRepository)
	}
	type args struct {
		id  domain.QuestionID
		req QuestionRequest
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When updating an existing question, then save the changes",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(defaultQuestion, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Cond(func(q domain.Question) bool {
						answers := q.Answers()

						return q.Description() == "What is the capital of Switzerland?" &&
							len(answers) == 1 &&
							answers[0].Description() == "Bern"
					})).Return(nil)
				},
			},

			args: args{
				id: defaultQuestion.ID(),
				req: QuestionRequest{
					Description: "What is the capital of Switzerland?",
					Visibility:  domain.Private,
					Answer: []AnswerRequest{
						{Description: "Bern", Position: 0, Correct: true},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "When the question ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(domain.Question{}, ErrEntityNotFound)
				},
			},
			args: args{
				id: defaultQuestion.ID(),
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(domain.Question{}, errors.New("database is down"))
				},
			},
			args: args{
				id: defaultQuestion.ID(),
			},
			wantErr: ErrPersistence,
		},
		{
			name: "When the description is empty, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(defaultQuestion, nil)
				},
			},
			args: args{
				id: defaultQuestion.ID(),
				req: QuestionRequest{
					Description: "",
					Visibility:  domain.Private,
					Answer:      nil,
				},
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When an answer position is negative, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(defaultQuestion, nil)
				},
			},
			args: args{
				id: defaultQuestion.ID(),
				req: QuestionRequest{
					Description: "What is the capital of Switzerland?",
					Visibility:  domain.Private,
					Answer: []AnswerRequest{
						{Description: "Bern", Position: -1, Correct: true},
					},
				},
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When saving the question fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), defaultQuestion.ID()).Return(defaultQuestion, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				id: defaultQuestion.ID(),
				req: QuestionRequest{
					Description: "What is the capital of Switzerland?",
					Visibility:  domain.Private,
					Answer:      nil,
				},
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockQuestionRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuestionService{repository: repo}
			err := s.Update(context.Background(), actor, tt.args.id, tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestionService_Delete(t *testing.T) {
	actor := domain.UserID(uuid.NewV7())
	id := domain.QuestionID(uuid.Nil())
	defaultQuestion, _ := domain.NewQuestion("What is the capital of Switzerland?", actor, domain.Private)

	type mocks struct {
		repository func(r *MockQuestionRepository)
	}
	type args struct {
		id domain.QuestionID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When the question exists, then delete it",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), id).Return(defaultQuestion, nil)
					r.EXPECT().Delete(gomock.Any(), id).Return(nil)
				},
			},

			args: args{
				id: id,
			},
			wantErr: nil,
		},
		{
			name: "When the question ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), id).Return(domain.Question{}, ErrEntityNotFound)
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
				repository: func(r *MockQuestionRepository) {
					r.EXPECT().Find(gomock.Any(), id).Return(defaultQuestion, nil)
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
			repo := NewMockQuestionRepository(ctrl)
			tt.mocks.repository(repo)

			s := &QuestionService{repository: repo}
			err := s.Delete(context.Background(), actor, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
