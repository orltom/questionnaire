package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"uuid"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type quizService interface {
	Create(ctx context.Context, title, description string) (domain.Quiz, error)
	Get(ctx context.Context, id domain.QuizID) (domain.Quiz, error)
	Update(ctx context.Context, id domain.QuizID, title, description string, visibility domain.Visibility) error
	Delete(ctx context.Context, id domain.QuizID) error
}

type quizHandler struct {
	service quizService
}

func (h quizHandler) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req api.CreateQuiz

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	resp, err := h.service.Create(ctx, req.Title, req.Description)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create quiz", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(toAPIQuiz(resp))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h quizHandler) GetQuiz(w http.ResponseWriter, r *http.Request, quizID api.QuizId) {
	ctx := r.Context()

	quiz, err := h.service.Get(ctx, domain.QuizID(quizID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to load quiz", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(toAPIQuiz(quiz))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h quizHandler) UpdateQuiz(w http.ResponseWriter, r *http.Request, quizID api.QuizId) {
	ctx := r.Context()

	var req api.UpdateQuiz

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	err = h.service.Update(ctx, domain.QuizID(quizID), req.Title, req.Description, domain.Visibility(req.Visibility))
	if err != nil {
		slog.ErrorContext(ctx, "failed to update quiz", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h quizHandler) DeleteQuiz(w http.ResponseWriter, r *http.Request, quizID api.QuizId) {
	ctx := r.Context()

	err := h.service.Delete(ctx, domain.QuizID(quizID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete quiz", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toAPIQuiz(resp domain.Quiz) api.Quiz {
	dst := make([]api.QuizQuestion, len(resp.Questions()))
	for i, q := range resp.Questions() {
		dst[i] = api.QuizQuestion{
			QuestionId: api.QuestionId(q.ID()),
			Position:   q.Position(),
		}
	}

	return api.Quiz{
		Id:          uuid.UUID(resp.ID()),
		Title:       resp.Title(),
		Description: resp.Description(),
		Visibility:  api.Visibility(resp.Visibility()),
		Questions:   dst,
	}
}
