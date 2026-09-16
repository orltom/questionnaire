package persistence

import (
	"context"
	"sync"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

var _ application.ParticipationRepository = (*InMemoryParticipationRepository)(nil)

type InMemoryParticipationRepository struct {
	mx    sync.RWMutex
	cache map[domain.ParticipationID]domain.Participation
}

func NewInMemoryParticipationRepository() *InMemoryParticipationRepository {
	return &InMemoryParticipationRepository{
		mx:    sync.RWMutex{},
		cache: make(map[domain.ParticipationID]domain.Participation),
	}
}

func (r *InMemoryParticipationRepository) Create(_ context.Context, participation domain.Participation) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.cache[participation.ID()] = participation

	return nil
}

func (r *InMemoryParticipationRepository) Update(_ context.Context, participation domain.Participation) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	_, ok := r.cache[participation.ID()]
	if !ok {
		return application.ErrEntityNotFound
	}

	r.cache[participation.ID()] = participation

	return nil
}

func (r *InMemoryParticipationRepository) FindByUser(
	_ context.Context,
	id domain.ChallengeID,
	user domain.UserID,
) (domain.Participation, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	for _, p := range r.cache {
		if p.ChallengeID() == id && p.UserID() == user {
			return p, nil
		}
	}

	return domain.Participation{}, application.ErrEntityNotFound
}

func (r *InMemoryParticipationRepository) FindByChallenge(
	_ context.Context,
	id domain.ChallengeID,
) ([]domain.Participation, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	participations := make([]domain.Participation, 0)

	for _, p := range r.cache {
		if p.ChallengeID() == id {
			participations = append(participations, p)
		}
	}

	return participations, nil
}
