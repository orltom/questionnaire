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

func TestQuestionHandler_CreateQuestion(t *testing.T) {
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?", actor, domain.Public)

	type fields struct {
		service func(s *MockquestionService)
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
		want       *api.Question
	}{
		{
			name: "When the request is valid, then create the question and return it",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Create(gomock.Any(), actor, application.QuestionRequest{
						Description: "What is the largest animal?",
						Visibility:  domain.Public,
						Labels:      []domain.Label{"animals"},
						Answer: []application.AnswerRequest{
							{Description: "Blue whale", Position: 1, Correct: true},
						},
					}).Return(defaultQuestion, nil)
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"description":"What is the largest animal?","visibility":"public","labels":["animals"],"answers":[{"description":"Blue whale","position":1,"correct":true}]}`,
			},
			wantStatus: http.StatusCreated,
			want: &api.Question{
				Id:          api.QuestionId(defaultQuestion.ID()),
				Owner:       api.QuestionId(actor),
				Description: "What is the largest animal?",
				Visibility:  api.Visibility(domain.Public),
				Labels:      []string{},
				Answers:     []api.Answer{},
			},
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				ctx:  func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				body: `{"description":"What is the largest animal?","visibility":"public","answers":[]}`,
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				ctx:  actorContext,
				body: `{"description":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the question is invalid, then return bad request",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Create(gomock.Any(), actor, gomock.Any()).Return(domain.Question{}, application.ErrInvalidArguments)
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"description":"","visibility":"public","answers":[]}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Create(gomock.Any(), actor, gomock.Any()).Return(domain.Question{}, errors.New("database is down"))
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"description":"What is the largest animal?","visibility":"public","answers":[]}`,
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
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodPost, "/questions", strings.NewReader(tt.args.body))

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
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?", actor, domain.Public)

	type fields struct {
		service func(s *MockquestionService)
	}
	type args struct {
		ctx        func(t *testing.T) context.Context
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
					s.EXPECT().Get(gomock.Any(), actor, defaultQuestion.ID()).Return(defaultQuestion, nil)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusOK,
			want: &api.Question{
				Id:          api.QuestionId(defaultQuestion.ID()),
				Owner:       api.QuestionId(actor),
				Description: "What is the largest animal?",
				Visibility:  api.Visibility(domain.Public),
				Labels:      []string{},
				Answers:     []api.Answer{},
			},
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				ctx:        func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the user is not allowed to view the question, then return forbidden",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Question{}, application.ErrForbidden)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the question does not exist, then return not found",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Question{}, application.ErrEntityNotFound)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Question{}, errors.New("database is down"))
				},
			},
			args: args{
				ctx:        actorContext,
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
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodGet, "/questions/"+tt.args.questionID.String(), nil)

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
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?", actor, domain.Public)

	type fields struct {
		service func(s *MockquestionService)
	}
	type args struct {
		ctx        func(t *testing.T) context.Context
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
					s.EXPECT().Update(gomock.Any(), actor, defaultQuestion.ID(), application.QuestionRequest{
						Description: "What is the fastest animal?",
						Visibility:  domain.Private,
						Labels:      []domain.Label{},
						Answer: []application.AnswerRequest{
							{Description: "Peregrine falcon", Position: 1, Correct: true},
						},
					}).Return(nil)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","visibility":"private","answers":[{"description":"Peregrine falcon","position":1,"correct":true}]}`,
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				ctx:        func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","visibility":"private","answers":[]}`,
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the user is not allowed to edit the question, then return forbidden",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Update(gomock.Any(), actor, gomock.Any(), gomock.Any()).Return(application.ErrForbidden)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","visibility":"public","answers":[]}`,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the question does not exist, then return not found",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Update(gomock.Any(), actor, gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","visibility":"public","answers":[]}`,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Update(gomock.Any(), actor, gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
				body:       `{"description":"What is the fastest animal?","visibility":"public","answers":[]}`,
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
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodPut, "/questions/"+tt.args.questionID.String(), strings.NewReader(tt.args.body))

			h.UpdateQuestion(rec, req, tt.args.questionID)

			if rec.Code != tt.wantStatus {
				t.Errorf("UpdateQuestion() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestQuestionHandler_DeleteQuestion(t *testing.T) {
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	defaultQuestion, _ := domain.NewQuestion("What is the largest animal?", actor, domain.Public)

	type fields struct {
		service func(s *MockquestionService)
	}
	type args struct {
		ctx        func(t *testing.T) context.Context
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
					s.EXPECT().Delete(gomock.Any(), actor, defaultQuestion.ID()).Return(nil)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockquestionService) {},
			},
			args: args{
				ctx:        func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the user is not allowed to delete the question, then return forbidden",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Delete(gomock.Any(), actor, gomock.Any()).Return(application.ErrForbidden)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the question does not exist, then return not found",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Delete(gomock.Any(), actor, gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				ctx:        actorContext,
				questionID: api.QuestionId(defaultQuestion.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockquestionService) {
					s.EXPECT().Delete(gomock.Any(), actor, gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				ctx:        actorContext,
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
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodDelete, "/questions/"+tt.args.questionID.String(), nil)

			h.DeleteQuestion(rec, req, tt.args.questionID)

			if rec.Code != tt.wantStatus {
				t.Errorf("DeleteQuestion() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
