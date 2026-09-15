package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"uuid"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type challengeService interface {
	Create(ctx context.Context, actor domain.UserID, req application.ChallengeRequest) (domain.Challenge, error)
	Get(ctx context.Context, actor domain.UserID, id domain.ChallengeID) (domain.Challenge, error)
	Close(ctx context.Context, actor domain.UserID, id domain.ChallengeID, now time.Time) error
	Delete(ctx context.Context, actor domain.UserID, id domain.ChallengeID) error
}

type participationService interface {
	Submit(
		ctx context.Context,
		actor domain.UserID,
		id domain.ChallengeID,
		questionID domain.QuestionID,
		answerID domain.AnswerID,
		now time.Time,
	) error
	Results(ctx context.Context, actor domain.UserID, id domain.ChallengeID) ([]application.Result, error)
}

type challengeHandler struct {
	service       challengeService
	participation participationService
	now           func() time.Time
}

func (h challengeHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	var dto api.CreateChallenge

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	invitees := make([]domain.UserID, 0)
	if dto.Invitees != nil {
		invitees = make([]domain.UserID, len(*dto.Invitees))
		for i, id := range *dto.Invitees {
			invitees[i] = domain.UserID(id)
		}
	}

	challenge, err := h.service.Create(ctx, actor, application.ChallengeRequest{
		QuizID:   domain.QuizID(dto.QuizId),
		StartsAt: dto.StartsAt,
		EndsAt:   dto.EndsAt,
		Access:   domain.Access(dto.Access),
		Invitees: invitees,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create challenge", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(toAPIChallenge(challenge, h.now()))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h challengeHandler) GetChallenge(w http.ResponseWriter, r *http.Request, challengeID api.ChallengeId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	challenge, err := h.service.Get(ctx, actor, domain.ChallengeID(challengeID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to load challenge", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(toAPIChallenge(challenge, h.now()))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h challengeHandler) DeleteChallenge(w http.ResponseWriter, r *http.Request, challengeID api.ChallengeId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	err = h.service.Delete(ctx, actor, domain.ChallengeID(challengeID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete challenge", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h challengeHandler) CloseChallenge(w http.ResponseWriter, r *http.Request, challengeID api.ChallengeId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	err = h.service.Close(ctx, actor, domain.ChallengeID(challengeID), h.now())
	if err != nil {
		slog.ErrorContext(ctx, "failed to close challenge", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h challengeHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request, challengeID api.ChallengeId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	var dto api.SubmitAnswer

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	err = h.participation.Submit(
		ctx,
		actor,
		domain.ChallengeID(challengeID),
		domain.QuestionID(dto.QuestionId),
		domain.AnswerID(dto.AnswerId),
		h.now(),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to submit answer", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h challengeHandler) GetChallengeResults(w http.ResponseWriter, r *http.Request, challengeID api.ChallengeId) {
	ctx := r.Context()

	actor, err := actorFrom(ctx)
	if err != nil {
		unauthorized(ctx, w, err)

		return
	}

	results, err := h.participation.Results(ctx, actor, domain.ChallengeID(challengeID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to load challenge results", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	dst := make([]api.Result, len(results))
	for i, res := range results {
		dst[i] = api.Result{
			UserId: uuid.UUID(res.User),
			Points: res.Points,
		}
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(dst)
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func toAPIChallenge(challenge domain.Challenge, now time.Time) api.Challenge {
	snapshot := challenge.Quiz()

	questions := make([]api.ChallengeQuestion, len(snapshot.Questions()))
	for i, q := range snapshot.Questions() {
		answers := make([]api.ChallengeAnswer, len(q.Answers()))
		for j, a := range q.Answers() {
			answers[j] = api.ChallengeAnswer{
				Id:          uuid.UUID(a.ID()),
				Description: a.Description(),
				Position:    a.Position(),
			}
		}

		questions[i] = api.ChallengeQuestion{
			Id:          uuid.UUID(q.ID()),
			Description: q.Description(),
			Position:    q.Position(),
			Answers:     answers,
		}
	}

	invitees := make([]uuid.UUID, len(challenge.Invitees()))
	for i, id := range challenge.Invitees() {
		invitees[i] = uuid.UUID(id)
	}

	return api.Challenge{
		Id:       uuid.UUID(challenge.ID()),
		Owner:    uuid.UUID(challenge.Owner()),
		StartsAt: challenge.StartsAt(),
		EndsAt:   challenge.EndsAt(),
		ClosedAt: challenge.ClosedAt(),
		Access:   api.Access(challenge.Access()),
		Invitees: invitees,
		State:    api.ChallengeState(challenge.State(now)),
		Quiz: api.ChallengeQuiz{
			QuizId:    uuid.UUID(snapshot.QuizID()),
			Title:     snapshot.Title(),
			Questions: questions,
		},
	}
}
