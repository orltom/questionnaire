package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestQuizHandler_CreateQuiz(t *testing.T) {
	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
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
					s.EXPECT().Create(gomock.Any(), "Animals", "Questions about animals").Return(defaultQuiz, nil)
				},
			},
			args: args{
				body: `{"title":"Animals","description":"Questions about animals","visibility":"public"}`,
			},
			wantStatus: http.StatusCreated,
			want: &api.Quiz{
				Id:          api.QuizId(defaultQuiz.ID()),
				Title:       "Animals",
				Description: "Questions about animals",
				Visibility:  api.Visibility(domain.Public),
				Questions:   []api.QuizQuestion{},
			},
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				body: `{"title":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the quiz is invalid, then return bad request",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Quiz{}, application.ErrInvalidArguments)
				},
			},
			args: args{
				body: `{"title":"","description":"","visibility":"public"}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Quiz{}, errors.New("database is down"))
				},
			},
			args: args{
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
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/quizzes", strings.NewReader(tt.args.body))

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
	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
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
					s.EXPECT().Get(gomock.Any(), defaultQuiz.ID()).Return(defaultQuiz, nil)
				},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusOK,
			want: &api.Quiz{
				Id:          api.QuizId(defaultQuiz.ID()),
				Title:       "Animals",
				Description: "Questions about animals",
				Visibility:  api.Visibility(domain.Public),
				Questions:   []api.QuizQuestion{},
			},
		},
		{
			name: "When the quiz does not exist, then return not found",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.Quiz{}, application.ErrEntityNotFound)
				},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.Quiz{}, errors.New("database is down"))
				},
			},
			args: args{
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
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/quizzes/"+tt.args.quizID.String(), nil)

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
	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
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
					s.EXPECT().Update(gomock.Any(), defaultQuiz.ID(), "Plants", "Questions about plants", domain.Private).Return(nil)
				},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"Plants","description":"Questions about plants","visibility":"private","questions":[]}`,
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquizService) {},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the quiz is invalid, then return bad request",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrInvalidArguments)
				},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"","description":"","visibility":"public","questions":[]}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the quiz does not exist, then return not found",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
				body:   `{"title":"Plants","description":"Questions about plants","visibility":"public","questions":[]}`,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
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
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/quizzes/"+tt.args.quizID.String(), strings.NewReader(tt.args.body))

			h.UpdateQuiz(rec, req, tt.args.quizID)

			if rec.Code != tt.wantStatus {
				t.Errorf("UpdateQuiz() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestQuizHandler_DeleteQuiz(t *testing.T) {
	defaultQuiz, _ := domain.NewQuiz("Animals", "Questions about animals", domain.Public)

	type fields struct {
		service func(s *MockquizService)
	}
	type args struct {
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
					s.EXPECT().Delete(gomock.Any(), defaultQuiz.ID()).Return(nil)
				},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the quiz does not exist, then return not found",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				quizID: api.QuizId(defaultQuiz.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquizService) {
					s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
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
			req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/quizzes/"+tt.args.quizID.String(), nil)

			h.DeleteQuiz(rec, req, tt.args.quizID)

			if rec.Code != tt.wantStatus {
				t.Errorf("DeleteQuiz() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
