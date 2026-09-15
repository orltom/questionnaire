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

type catalogService interface {
	Create(ctx context.Context, title, description string) (domain.Catalog, error)
	Get(ctx context.Context, id domain.CatalogID) (domain.Catalog, error)
	Update(ctx context.Context, id domain.CatalogID, title, description string) error
	Delete(ctx context.Context, id domain.CatalogID) error
}

type catalogHandler struct {
	service catalogService
}

func (h catalogHandler) CreateCatalog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto api.CreateCatalog

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	c, err := h.service.Create(ctx, dto.Title, dto.Description)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create catalog", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(toAPICatalog(c))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h catalogHandler) GetCatalog(w http.ResponseWriter, r *http.Request, catalogID api.CatalogId) {
	ctx := r.Context()

	catalog, err := h.service.Get(ctx, domain.CatalogID(catalogID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to load catalog", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(toAPICatalog(catalog))
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)

		return
	}
}

func (h catalogHandler) UpdateCatalog(w http.ResponseWriter, r *http.Request, catalogID api.CatalogId) {
	ctx := r.Context()

	var dto api.UpdateCatalog

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request body", "error", err)
		http.Error(w, "malformed request body", http.StatusBadRequest)

		return
	}

	err = h.service.Update(ctx, domain.CatalogID(catalogID), dto.Title, dto.Description)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update catalog", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h catalogHandler) DeleteCatalog(w http.ResponseWriter, r *http.Request, catalogID api.CatalogId) {
	ctx := r.Context()

	err := h.service.Delete(ctx, domain.CatalogID(catalogID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete catalog", "error", err)
		status, msg := httpError(err)
		http.Error(w, msg, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toAPICatalog(c domain.Catalog) api.Catalog {
	ids := make([]uuid.UUID, len(c.QuestionIDs()))
	for i, id := range c.QuestionIDs() {
		ids[i] = uuid.UUID(id)
	}

	return api.Catalog{
		Id:          api.CatalogId(c.ID()),
		Title:       c.Title(),
		Description: c.Description(),
		QuestionIds: ids,
	}
}
