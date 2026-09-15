package database

import (
	"context"
	"sync"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

var _ application.QuestionRepository = (*InMemoryRepository[domain.QuestionID, domain.Question])(nil)

func NewInMemoryQuestionRepository() *InMemoryRepository[domain.QuestionID, domain.Question] {
	return &InMemoryRepository[domain.QuestionID, domain.Question]{
		cache: make(map[domain.QuestionID]domain.Question),
		mx:    sync.RWMutex{},
		id: func(question domain.Question) domain.QuestionID {
			return question.ID()
		},
	}
}

var _ application.QuizRepository = (*InMemoryRepository[domain.QuizID, domain.Quiz])(nil)

func NewInMemoryQuizRepository() *InMemoryRepository[domain.QuizID, domain.Quiz] {
	return &InMemoryRepository[domain.QuizID, domain.Quiz]{
		cache: make(map[domain.QuizID]domain.Quiz),
		mx:    sync.RWMutex{},
		id: func(quiz domain.Quiz) domain.QuizID {
			return quiz.ID()
		},
	}
}

var _ application.ChallengeRepository = (*InMemoryRepository[domain.ChallengeID, domain.Challenge])(nil)

func NewInMemoryChallengeRepository() *InMemoryRepository[domain.ChallengeID, domain.Challenge] {
	return &InMemoryRepository[domain.ChallengeID, domain.Challenge]{
		cache: make(map[domain.ChallengeID]domain.Challenge),
		mx:    sync.RWMutex{},
		id: func(challenge domain.Challenge) domain.ChallengeID {
			return challenge.ID()
		},
	}
}

type InMemoryRepository[K comparable, V any] struct {
	mx    sync.RWMutex
	cache map[K]V
	id    func(V) K
}

func (i *InMemoryRepository[K, V]) Create(_ context.Context, entity V) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	i.cache[i.id(entity)] = entity

	return nil
}

func (i *InMemoryRepository[K, V]) Find(_ context.Context, id K) (V, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()

	c, ok := i.cache[id]
	if !ok {
		return c, application.ErrEntityNotFound
	}

	return c, nil
}

func (i *InMemoryRepository[K, V]) Update(_ context.Context, entity V) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	_, ok := i.cache[i.id(entity)]
	if !ok {
		return application.ErrEntityNotFound
	}

	i.cache[i.id(entity)] = entity

	return nil
}

func (i *InMemoryRepository[K, V]) Delete(_ context.Context, id K) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	_, ok := i.cache[id]
	if !ok {
		return application.ErrEntityNotFound
	}

	delete(i.cache, id)

	return nil
}

func (i *InMemoryRepository[K, V]) Exists(_ context.Context, id K) (bool, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()

	_, ok := i.cache[id]

	return ok, nil
}
