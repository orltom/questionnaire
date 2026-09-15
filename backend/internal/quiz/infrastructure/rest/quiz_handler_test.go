package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"uuid"

	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/identity"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestQuizHandler_CreateQuiz(t *testing.T) {
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", actor, domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
		ctx  func(t *testing.T) context.Context
		body string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       *api.Quiz
	}{
		{
			name: "When the request is valid, then create the quiz and return it",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Create(gomock.Any(), actor, "Animals", "Questions about animals", domain.Public).Return(defaultQuiz, nil)
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"title":"Animals","description":"Questions about animals","visibility":"public"}`,
			},
			wantStatus: http.StatusCreated,
			want: &api.Quiz{
				Id:          api.QuizId(defaultQuiz.ID()),
				Owner:       uuid.UUID(actor),
				Title:       "Animals",
				Description: "Questions about animals",
				Visibility:  api.Visibility(domain.Public),
				Questions:   []api.QuizQuestion{},
			},
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				ctx:  func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				body: `{"title":"Animals","description":"Questions about animals","visibility":"public"}`,
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				ctx:  actorContext,
				body: `{"title":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the quiz is invalid, then return bad request",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Create(gomock.Any(), actor, gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Quiz{}, application.ErrInvalidArguments)
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"title":"","description":"","visibility":"public"}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Create(gomock.Any(), actor, gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Quiz{}, errors.New("database is down"))
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"title":"Animals","description":"Questions about animals","visibility":"public"}`,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquizService(ctrl)
			tt.fields.service(service)

			h := quizHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodPost, "/quizzes", strings.NewReader(tt.args.body))

			h.CreateQuiz(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("CreateQuiz() status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.want == nil {
				return
			}

			var got api.Quiz
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("CreateQuiz() response body is not a quiz: %v", err)
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("CreateQuiz() = %+v, want %+v", got, *tt.want)
			}
		})
	}
}

func TestQuizHandler_GetQuiz(t *testing.T) {
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", actor, domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
		ctx    func(t *testing.T) context.Context
		quizID api.QuizId
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       *api.Quiz
	}{
		{
			name: "When the quiz exists, then return it",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Get(gomock.Any(), actor, defaultQuiz.ID()).Return(defaultQuiz, nil)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusOK,
			want: &api.Quiz{
				Id:          api.QuizId(defaultQuiz.ID()),
				Owner:       uuid.UUID(actor),
				Title:       "Animals",
				Description: "Questions about animals",
				Visibility:  api.Visibility(domain.Public),
				Questions:   []api.QuizQuestion{},
			},
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				ctx:    func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the user is not allowed to view the quiz, then return forbidden",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Quiz{}, application.ErrForbidden)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the quiz does not exist, then return not found",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Quiz{}, application.ErrEntityNotFound)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Quiz{}, errors.New("database is down"))
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquizService(ctrl)
			tt.fields.service(service)

			h := quizHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodGet, "/quizzes/"+tt.args.quizID.String(), nil)

			h.GetQuiz(rec, req, tt.args.quizID)

			if rec.Code != tt.wantStatus {
				t.Fatalf("GetQuiz() status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.want == nil {
				return
			}

			var got api.Quiz
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("GetQuiz() response body is not a quiz: %v", err)
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("GetQuiz() = %+v, want %+v", got, *tt.want)
			}
		})
	}
}

func TestQuizHandler_UpdateQuiz(t *testing.T) {
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", actor, domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
		ctx    func(t *testing.T) context.Context
		quizID api.QuizId
		body   string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "When the request is valid, then update the quiz",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Update(gomock.Any(), actor, defaultQuiz.ID(), "Plants", "Questions about plants", domain.Private).Return(nil)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"Plants","description":"Questions about plants","visibility":"private","questions":[]}`,
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				ctx:    func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"Plants","description":"Questions about plants","visibility":"private","questions":[]}`,
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the user is not allowed to edit the quiz, then return forbidden",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Update(gomock.Any(), actor, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrForbidden)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"Plants","description":"Questions about plants","visibility":"public","questions":[]}`,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the quiz does not exist, then return not found",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Update(gomock.Any(), actor, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"Plants","description":"Questions about plants","visibility":"public","questions":[]}`,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Update(gomock.Any(), actor, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"Plants","description":"Questions about plants","visibility":"public","questions":[]}`,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquizService(ctrl)
			tt.fields.service(service)

			h := quizHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodPut, "/quizzes/"+tt.args.quizID.String(), strings.NewReader(tt.args.body))

			h.UpdateQuiz(rec, req, tt.args.quizID)

			if rec.Code != tt.wantStatus {
				t.Errorf("UpdateQuiz() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestQuizHandler_DeleteQuiz(t *testing.T) {
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", actor, domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
		ctx    func(t *testing.T) context.Context
		quizID api.QuizId
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "When the quiz exists, then delete it",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Delete(gomock.Any(), actor, defaultQuiz.ID()).Return(nil)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				ctx:    func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the user is not allowed to delete the quiz, then return forbidden",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Delete(gomock.Any(), actor, gomock.Any()).Return(application.ErrForbidden)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the quiz does not exist, then return not found",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Delete(gomock.Any(), actor, gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Delete(gomock.Any(), actor, gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				ctx:    actorContext,
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquizService(ctrl)
			tt.fields.service(service)

			h := quizHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodDelete, "/quizzes/"+tt.args.quizID.String(), nil)

			h.DeleteQuiz(rec, req, tt.args.quizID)

			if rec.Code != tt.wantStatus {
				t.Errorf("DeleteQuiz() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
