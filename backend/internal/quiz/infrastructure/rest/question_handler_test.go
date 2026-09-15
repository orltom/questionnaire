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

func TestQuestionHandler_CreateQuestion(t *testing.T) {
	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?")

	type fields struct {
		service func(s *MockquestionService)
	}
	type args struct {
		body string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       *api.Question
	}{
		{
			name: "When the request is valid, then create the question and return it",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Create(gomock.Any(), application.QuestionRequest{
						Description: "What is the largest animal?",
						Answer: []application.AnswerRequest{
							{Description: "Blue whale", Position: 1, Count: 0},
						},
					}).Return(defaultQuestion, nil)
				},
			},
			args: args{
				body: `{"description":"What is the largest animal?","answers":[{"description":"Blue whale","position":1,"count":0}]}`,
			},
			wantStatus: http.StatusCreated,
			want: &api.Question{
				Id:          api.QuestionId(defaultQuestion.ID()),
				Description: "What is the largest animal?",
				Answers:     []api.Answer{},
			},
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				body: `{"description":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the question is invalid, then return bad request",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.Question{}, application.ErrInvalidArguments)
				},
			},
			args: args{
				body: `{"description":"","answers":[]}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.Question{}, errors.New("database is down"))
				},
			},
			args: args{
				body: `{"description":"What is the largest animal?","answers":[]}`,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquestionService(ctrl)
			tt.fields.service(service)

			h := questionHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/questions", strings.NewReader(tt.args.body))

			h.CreateQuestion(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("CreateQuestion() status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.want == nil {
				return
			}

			var got api.Question
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("CreateQuestion() response body is not a question: %v", err)
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("CreateQuestion() = %+v, want %+v", got, *tt.want)
			}
		})
	}
}

func TestQuestionHandler_GetQuestion(t *testing.T) {
	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?")

	type fields struct {
		service func(s *MockquestionService)
	}
	type args struct {
		questionID api.QuestionId
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       *api.Question
	}{
		{
			name: "When the question exists, then return it",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Get(gomock.Any(), defaultQuestion.ID()).Return(defaultQuestion, nil)
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusOK,
			want: &api.Question{
				Id:          api.QuestionId(defaultQuestion.ID()),
				Description: "What is the largest animal?",
				Answers:     []api.Answer{},
			},
		},
		{
			name: "When the question does not exist, then return not found",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.Question{}, application.ErrEntityNotFound)
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.Question{}, errors.New("database is down"))
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquestionService(ctrl)
			tt.fields.service(service)

			h := questionHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/questions/"+tt.args.questionID.String(), nil)

			h.GetQuestion(rec, req, tt.args.questionID)

			if rec.Code != tt.wantStatus {
				t.Fatalf("GetQuestion() status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.want == nil {
				return
			}

			var got api.Question
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("GetQuestion() response body is not a question: %v", err)
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("GetQuestion() = %+v, want %+v", got, *tt.want)
			}
		})
	}
}

func TestQuestionHandler_UpdateQuestion(t *testing.T) {
	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?")

	type fields struct {
		service func(s *MockquestionService)
	}
	type args struct {
		questionID api.QuestionId
		body       string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "When the request is valid, then update the question",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Update(gomock.Any(), defaultQuestion.ID(), application.QuestionRequest{
						Description: "What is the fastest animal?",
						Answer: []application.AnswerRequest{
							{Description: "Peregrine falcon", Position: 1, Count: 0},
						},
					}).Return(nil)
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","answers":[{"description":"Peregrine falcon","position":1,"count":0}]}`,
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the question is invalid, then return bad request",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrInvalidArguments)
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"","answers":[]}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the question does not exist, then return not found",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","answers":[]}`,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","answers":[]}`,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquestionService(ctrl)
			tt.fields.service(service)

			h := questionHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/questions/"+tt.args.questionID.String(), strings.NewReader(tt.args.body))

			h.UpdateQuestion(rec, req, tt.args.questionID)

			if rec.Code != tt.wantStatus {
				t.Errorf("UpdateQuestion() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestQuestionHandler_DeleteQuestion(t *testing.T) {
	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?")

	type fields struct {
		service func(s *MockquestionService)
	}
	type args struct {
		questionID api.QuestionId
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "When the question exists, then delete it",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Delete(gomock.Any(), defaultQuestion.ID()).Return(nil)
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the question does not exist, then return not found",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockquestionService(ctrl)
			tt.fields.service(service)

			h := questionHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/questions/"+tt.args.questionID.String(), nil)

			h.DeleteQuestion(rec, req, tt.args.questionID)

			if rec.Code != tt.wantStatus {
				t.Errorf("DeleteQuestion() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
