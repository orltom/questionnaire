package persistence

import (
	"context"
	"sync"

	"gitlab.com/orltom/questionnaire/backend/internal/identity"
)

var _ identity.UserRepository = (*InMemoryUserRepository)(nil)

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		cache: map[issuerSubject]identity.User{},
		mutex: sync.RWMutex{},
	}
}

type issuerSubject struct {
	issuer  string
	subject string
}

type InMemoryUserRepository struct {
	cache map[issuerSubject]identity.User
	mutex sync.RWMutex
}

func (i *InMemoryUserRepository) Create(_ context.Context, user identity.User) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.cache[issuerSubject{user.Issuer(), user.Subject()}] = user
	return nil
}

func (i *InMemoryUserRepository) FindBySubject(_ context.Context, issuer string, subject string) (identity.User, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	user, ok := i.cache[issuerSubject{issuer, subject}]
	if !ok {
		return user, identity.ErrEntityNotFound
	}
	return user, nil
}
