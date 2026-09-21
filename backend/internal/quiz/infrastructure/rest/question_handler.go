package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type questionService interface {
	Create(ctx context.Context, actor domain.UserID, req application.QuestionRequest) (domain.Question, error)
	Get(ctx context.Context, actor domain.UserID, id domain.QuestionID) (domain.Question, error)
	Update(ctx context.Context, actor domain.UserID, id domain.QuestionID, req application.QuestionRequest) error
	Delete(ctx context.Context, actor domain.UserID, id domain.QuestionID) error
}

type questionHandler struct {
	service questionService
}

func (h questionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	var dto api.CreateQuestion

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	q, err := h.service.Create(ctx, actor, toQuestionRequest(dto.Description, dto.Visibility, dto.Labels, dto.Answers))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create question", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(toAPIQuestion(q))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h questionHandler) GetQuestion(w http.ResponseWriter, r *http.Request, questionID api.QuestionId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	question, err := h.service.Get(ctx, actor, domain.QuestionID(questionID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to load question", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(toAPIQuestion(question))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h questionHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request, questionID api.QuestionId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	var dto api.UpdateQuestion

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	req := toQuestionRequest(dto.Description, dto.Visibility, dto.Labels, dto.Answers)

	err = h.service.Update(ctx, actor, domain.QuestionID(questionID), req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update question", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h questionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request, questionID api.QuestionId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	err = h.service.Delete(ctx, actor, domain.QuestionID(questionID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete question", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toQuestionRequest(description string, visibility api.Visibility, labels *[]string, answers []api.Answer) application.QuestionRequest {
	req := application.QuestionRequest{
		Description: description,
		Visibility:  domain.Visibility(visibility),
		Labels:      []domain.Label{},
		Answer:      make([]application.AnswerRequest, len(answers)),
	}

	if labels != nil {
		req.Labels = make([]domain.Label, len(*labels))
		for i, l := range *labels {
			req.Labels[i] = domain.Label(l)
		}
	}

	for i, a := range answers {
		req.Answer[i] = application.AnswerRequest{
			Description: a.Description,
			Position:    a.Position,
			Correct:     a.Correct,
		}
	}

	return req
}

func toAPIQuestion(q domain.Question) api.Question {
	answers := make([]api.Answer, len(q.Answers()))
	for i, a := range q.Answers() {
		answers[i] = api.Answer{
			Id:          api.QuestionId(a.ID()),
			Description: a.Description(),
			Position:    a.Position(),
			Correct:     a.Correct(),
		}
	}

	labels := make([]string, len(q.Labels()))
	for i, l := range q.Labels() {
		labels[i] = string(l)
	}

	return api.Question{
		Id:          api.QuestionId(q.ID()),
		Owner:       api.QuestionId(q.Owner()),
		Description: q.Description(),
		Visibility:  api.Visibility(q.Visibility()),
		Labels:      labels,
		Answers:     answers,
	}
}
