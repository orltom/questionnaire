package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"uuid"

	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/identity"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestChallengeHandler_GetChallenge(t *testing.T) {
	frozen := time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC)
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	answer := domain.NewSnapshotAnswer(domain.AnswerID(uuid.NewV7()), "Blue whale", 0, true)
	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(domain.QuestionID(uuid.NewV7()), "Largest animal?", 0, answer),
	)
	challenge, _ := domain.NewChallenge(actor, snapshot, frozen.Add(-time.Hour), nil, domain.AccessPublic)

	type fields struct {
		service func(s *MockchallengeService)
	}
	type args struct {
		ctx         func(t *testing.T) context.Context
		challengeID api.ChallengeId
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantState  api.ChallengeState
	}{
		{
			name: "When the challenge exists, then return it without the correct answers",
			fields: fields{
				service: func(s *MockchallengeService) {
					s.EXPECT().Get(gomock.Any(), actor, challenge.ID()).Return(challenge, nil)
				},
			},
			args: args{
				ctx:         actorContext,
				challengeID: api.ChallengeId(challenge.ID()),
			},
			wantStatus: http.StatusOK,
			wantState:  api.ChallengeState(domain.Running),
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockchallengeService) {},
			},
			args: args{
				ctx:         func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				challengeID: api.ChallengeId(challenge.ID()),
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the user is not invited, then return forbidden",
			fields: fields{
				service: func(s *MockchallengeService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Challenge{}, application.ErrForbidden)
				},
			},
			args: args{
				ctx:         actorContext,
				challengeID: api.ChallengeId(challenge.ID()),
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the challenge does not exist, then return not found",
			fields: fields{
				service: func(s *MockchallengeService) {
					s.EXPECT().Get(gomock.Any(), actor, gomock.Any()).Return(domain.Challenge{}, application.ErrEntityNotFound)
				},
			},
			args: args{
				ctx:         actorContext,
				challengeID: api.ChallengeId(challenge.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockchallengeService(ctrl)
			tt.fields.service(service)

			h := challengeHandler{
				service:       service,
				participation: NewMockparticipationService(ctrl),
				now:           func() time.Time { return frozen },
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodGet, "/challenges/"+tt.args.challengeID.String(), nil)

			h.GetChallenge(rec, req, tt.args.challengeID)

			if rec.Code != tt.wantStatus {
				t.Fatalf("GetChallenge() status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus != http.StatusOK {
				return
			}

			body := rec.Body.String()
			if strings.Contains(body, "correct") {
				t.Errorf("GetChallenge() leaked the answer key: %s", body)
			}

			var got api.ChallengeOverview
			if err := json.Unmarshal([]byte(body), &got); err != nil {
				t.Fatalf("GetChallenge() response body is not a challenge overview: %v", err)
			}
			if got.State != tt.wantState {
				t.Errorf("GetChallenge() state = %q, want %q", got.State, tt.wantState)
			}
			if len(got.Questions) != 1 {
				t.Errorf("GetChallenge() questions = %+v, want one question", got.Questions)
			}
		})
	}
}

func TestChallengeHandler_CreateChallenge(t *testing.T) {
	frozen := time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC)
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	answer := domain.NewSnapshotAnswer(domain.AnswerID(uuid.NewV7()), "Blue whale", 0, true)
	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(domain.QuestionID(uuid.NewV7()), "Largest animal?", 0, answer),
	)
	challenge, _ := domain.NewChallenge(actor, snapshot, frozen.Add(-time.Hour), nil, domain.AccessPublic)

	type fields struct {
		service func(s *MockchallengeService)
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
	}{
		{
			name: "When the request is valid, then create the challenge",
			fields: fields{
				service: func(s *MockchallengeService) {
					s.EXPECT().Create(gomock.Any(), actor, gomock.Cond(func(req application.ChallengeRequest) bool {
						return req.Access == domain.AccessPublic && req.EndsAt == nil
					})).Return(challenge, nil)
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"quizId":"` + uuid.NewV7().String() + `","startsAt":"2026-09-15T12:00:00Z","access":"public"}`,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				service: func(s *MockchallengeService) {},
			},
			args: args{
				ctx:  func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				body: `{"quizId":"` + uuid.NewV7().String() + `","startsAt":"2026-09-15T12:00:00Z","access":"public"}`,
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockchallengeService) {},
			},
			args: args{
				ctx:  actorContext,
				body: `{"quizId":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the quiz can not be viewed, then return forbidden",
			fields: fields{
				service: func(s *MockchallengeService) {
					s.EXPECT().Create(gomock.Any(), actor, gomock.Any()).Return(domain.Challenge{}, application.ErrForbidden)
				},
			},
			args: args{
				ctx:  actorContext,
				body: `{"quizId":"` + uuid.NewV7().String() + `","startsAt":"2026-09-15T12:00:00Z","access":"public"}`,
			},
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockchallengeService(ctrl)
			tt.fields.service(service)

			h := challengeHandler{
				service:       service,
				participation: NewMockparticipationService(ctrl),
				now:           func() time.Time { return frozen },
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodPost, "/challenges", strings.NewReader(tt.args.body))

			h.CreateChallenge(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("CreateChallenge() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestChallengeHandler_SubmitAnswer(t *testing.T) {
	frozen := time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC)
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	answer := domain.NewSnapshotAnswer(domain.AnswerID(uuid.NewV7()), "Blue whale", 0, true)
	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(domain.QuestionID(uuid.NewV7()), "Largest animal?", 0, answer),
	)
	challenge, _ := domain.NewChallenge(actor, snapshot, frozen.Add(-time.Hour), nil, domain.AccessPublic)
	questionID := challenge.Quiz().Questions()[0].ID()
	answerID := challenge.Quiz().Questions()[0].Answers()[0].ID()
	body := `{"questionId":"` + uuid.UUID(questionID).String() + `","answerId":"` + uuid.UUID(answerID).String() + `"}`

	type fields struct {
		participation func(s *MockparticipationService)
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
	}{
		{
			name: "When the answer is valid, then record it",
			fields: fields{
				participation: func(s *MockparticipationService) {
					s.EXPECT().Submit(gomock.Any(), actor, challenge.ID(), questionID, answerID, frozen).Return(nil)
				},
			},
			args: args{
				ctx:  actorContext,
				body: body,
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the user is not authenticated, then return unauthorized",
			fields: fields{
				participation: func(s *MockparticipationService) {},
			},
			args: args{
				ctx:  func(t *testing.T) context.Context { t.Helper(); return t.Context() },
				body: body,
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "When the challenge is not running, then return forbidden",
			fields: fields{
				participation: func(s *MockparticipationService) {
					s.EXPECT().Submit(gomock.Any(), actor, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrForbidden)
				},
			},
			args: args{
				ctx:  actorContext,
				body: body,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "When the answer is not part of the quiz, then return bad request",
			fields: fields{
				participation: func(s *MockparticipationService) {
					s.EXPECT().Submit(gomock.Any(), actor, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrInvalidArguments)
				},
			},
			args: args{
				ctx:  actorContext,
				body: body,
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			participation := NewMockparticipationService(ctrl)
			tt.fields.participation(participation)

			h := challengeHandler{
				service:       NewMockchallengeService(ctrl),
				participation: participation,
				now:           func() time.Time { return frozen },
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(tt.args.ctx(t), http.MethodPost, "/challenges/"+uuid.UUID(challenge.ID()).String()+"/answers", strings.NewReader(tt.args.body))

			h.SubmitAnswer(rec, req, api.ChallengeId(challenge.ID()))

			if rec.Code != tt.wantStatus {
				t.Errorf("SubmitAnswer() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestChallengeHandler_GetChallengeResults(t *testing.T) {
	frozen := time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC)
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	answer := domain.NewSnapshotAnswer(domain.AnswerID(uuid.NewV7()), "Blue whale", 0, true)
	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(domain.QuestionID(uuid.NewV7()), "Largest animal?", 0, answer),
	)
	challenge, _ := domain.NewChallenge(actor, snapshot, frozen.Add(-time.Hour), nil, domain.AccessPublic)

	ctrl := gomock.NewController(t)
	participation := NewMockparticipationService(ctrl)
	participation.EXPECT().Results(gomock.Any(), actor, challenge.ID()).
		Return([]application.Result{{User: actor, Points: 3}}, nil)

	h := challengeHandler{
		service:       NewMockchallengeService(ctrl),
		participation: participation,
		now:           func() time.Time { return frozen },
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(actorContext(t), http.MethodGet, "/challenges/x/results", nil)

	h.GetChallengeResults(rec, req, api.ChallengeId(challenge.ID()))

	if rec.Code != http.StatusOK {
		t.Fatalf("GetChallengeResults() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []api.Result
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("GetChallengeResults() response body is not a result list: %v", err)
	}
	if len(got) != 1 || got[0].Points != 3 || got[0].UserId != api.QuizId(actor) {
		t.Errorf("GetChallengeResults() = %+v, want one result with 3 points", got)
	}
}

func TestChallengeHandler_CloseChallenge(t *testing.T) {
	frozen := time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC)
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	answer := domain.NewSnapshotAnswer(domain.AnswerID(uuid.NewV7()), "Blue whale", 0, true)
	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(domain.QuestionID(uuid.NewV7()), "Largest animal?", 0, answer),
	)
	challenge, _ := domain.NewChallenge(actor, snapshot, frozen.Add(-time.Hour), nil, domain.AccessPublic)

	ctrl := gomock.NewController(t)
	service := NewMockchallengeService(ctrl)
	service.EXPECT().Close(gomock.Any(), actor, challenge.ID(), frozen).Return(nil)

	h := challengeHandler{
		service:       service,
		participation: NewMockparticipationService(ctrl),
		now:           func() time.Time { return frozen },
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(actorContext(t), http.MethodPost, "/challenges/x/close", nil)

	h.CloseChallenge(rec, req, api.ChallengeId(challenge.ID()))

	if rec.Code != http.StatusNoContent {
		t.Errorf("CloseChallenge() status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestChallengeHandler_DeleteChallenge(t *testing.T) {
	frozen := time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC)
	session := identity.UserID(uuid.NewV7())
	actor := domain.UserID(session)
	actorContext := func(t *testing.T) context.Context {
		t.Helper()

		return WithActor(t.Context(), session)
	}

	answer := domain.NewSnapshotAnswer(domain.AnswerID(uuid.NewV7()), "Blue whale", 0, true)
	snapshot, _ := domain.NewQuizSnapshot(
		domain.QuizID(uuid.NewV7()),
		"Animals",
		domain.NewSnapshotQuestion(domain.QuestionID(uuid.NewV7()), "Largest animal?", 0, answer),
	)
	challenge, _ := domain.NewChallenge(actor, snapshot, frozen.Add(-time.Hour), nil, domain.AccessPublic)

	ctrl := gomock.NewController(t)
	service := NewMockchallengeService(ctrl)
	service.EXPECT().Delete(gomock.Any(), actor, challenge.ID()).Return(errors.New("database is down"))

	h := challengeHandler{
		service:       service,
		participation: NewMockparticipationService(ctrl),
		now:           func() time.Time { return frozen },
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(actorContext(t), http.MethodDelete, "/challenges/x", nil)

	h.DeleteChallenge(rec, req, api.ChallengeId(challenge.ID()))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("DeleteChallenge() status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
