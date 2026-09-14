package application

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

type CatalogRepository interface {
	Create(ctx context.Context, catalog domain.Catalog) error
	Find(ctx context.Context, id domain.CatalogID) (domain.Catalog, error)
	Update(ctx context.Context, catalog domain.Catalog) error
	Delete(ctx context.Context, id domain.CatalogID) error
}

type CatalogService struct {
	repository    CatalogRepository
	lookupService QuestionLookupService
}

func NewCatalogService(repository CatalogRepository, lookupService QuestionLookupService) *CatalogService {
	return &CatalogService{
		repository:    repository,
		lookupService: lookupService,
	}
}

func (s *CatalogService) Create(ctx context.Context, title, description string) (domain.Catalog, error) {
	catalog, err := domain.NewCatalog(title, description)
	if err != nil {
		return catalog, fmt.Errorf("%w: create catalog: %w", ErrInvalidArguments, err)
	}

	err = s.repository.Create(ctx, catalog)
	if err != nil {
		return catalog, fmt.Errorf("%w: save catalog: %w", ErrPersistence, err)
	}

	return catalog, nil
}

func (s *CatalogService) Get(ctx context.Context, id domain.CatalogID) (domain.Catalog, error) {
	catalog, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return domain.Catalog{}, ErrEntityNotFound
		}

		return catalog, fmt.Errorf("%w: load catalog: %w", ErrPersistence, err)
	}

	return catalog, nil
}

func (s *CatalogService) Update(ctx context.Context, id domain.CatalogID, title, description string) error {
	catalog, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: find catalog: %w", ErrPersistence, err)
	}

	err = catalog.Rename(title)
	if err != nil {
		return fmt.Errorf("%w: update catalog title: %w", ErrInvalidArguments, err)
	}

	err = catalog.ChangeDescription(description)
	if err != nil {
		return fmt.Errorf("%w: update catalog description: %w", ErrInvalidArguments, err)
	}

	err = s.repository.Update(ctx, catalog)
	if err != nil {
		return fmt.Errorf("%w: save catalog: %w", ErrPersistence, err)
	}

	return nil
}

func (s *CatalogService) Delete(ctx context.Context, id domain.CatalogID) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("%w: delete catalog: %w", ErrPersistence, err)
	}

	return nil
}

func (s *CatalogService) AddQuestions(ctx context.Context, id domain.CatalogID, ids ...domain.QuestionID) error {
	catalog, err := s.repository.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: could not find catalog with ID '%s': %w", ErrPersistence, id, err)
	}

	for _, qID := range ids {
		ok, err := s.lookupService.Exists(ctx, qID)
		if err != nil {
			return fmt.Errorf("%w: lookup question: %w", ErrPersistence, err)
		}

		if !ok {
			return fmt.Errorf("%w: question with ID '%s' does not exist", ErrInvalidArguments, qID)
		}

		err = catalog.AddQuestion(qID)
		if err != nil {
			return fmt.Errorf("%w: add question to catalog: %w", ErrInvalidArguments, err)
		}
	}

	err = s.repository.Update(ctx, catalog)
	if err != nil {
		return fmt.Errorf("%w: save catalog: %w", ErrPersistence, err)
	}

	return nil
}
