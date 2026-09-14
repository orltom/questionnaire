package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"uuid"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type questionHandler struct {
	service *application.QuestionService
}

func (h questionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto api.CreateQuestion

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	answers := make([]application.AnswerRequest, len(dto.Answers))
	for i, a := range dto.Answers {
		answers[i] = application.AnswerRequest{
			Description: a.Description,
			Position:    a.Position,
			Count:       a.Count,
		}
	}

	reqCreate := application.QuestionRequest{
		Description: dto.Description,
		Answer:      answers,
	}

	q, err := h.service.Create(ctx, reqCreate)
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

	question, err := h.service.Get(ctx, domain.QuestionID(questionID))
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

	var dto api.UpdateQuestion

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	answers := make([]application.AnswerRequest, len(dto.Answers))
	for i, a := range dto.Answers {
		answers[i] = application.AnswerRequest{
			Description: a.Description,
			Position:    a.Position,
			Count:       a.Count,
		}
	}

	reqUpdate := application.QuestionRequest{
		Description: dto.Description,
		Answer:      answers,
	}

	err = h.service.Update(ctx, domain.QuestionID(questionID), reqUpdate)
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

	err := h.service.Delete(ctx, domain.QuestionID(questionID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete question", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toAPIQuestion(q domain.Question) api.Question {
	answers := make([]api.Answer, len(q.Answers()))
	for i, a := range q.Answers() {
		answers[i] = api.Answer{
			Description: a.Description(),
			Count:       a.Count(),
			Position:    a.Position(),
		}
	}

	return api.Question{
		Id:          uuid.UUID(q.ID()),
		Description: q.Description(),
		Answers:     answers,
	}
}
