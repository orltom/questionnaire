package identity

import (
	"context"
	"errors"
	"fmt"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	Find(ctx context.Context, issuer string, subject string) (User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Authenticate(ctx context.Context, issuer string, subject string, name string) (User, error) {
	user, err := s.repository.Find(ctx, issuer, subject)
	if err != nil && !errors.Is(err, ErrEntityNotFound) {
		return User{}, fmt.Errorf("%w: load user: %w", ErrPersistence, err)
	}

	// first attempt to use application, create and save user...
	if errors.Is(err, ErrEntityNotFound) {
		user, err = NewUser(issuer, subject, name)
		if err != nil {
			return User{}, fmt.Errorf("%w: create user: %w", ErrInvalidArguments, err)
		}

		err = s.repository.Create(ctx, user)
		if err != nil {
			return User{}, fmt.Errorf("%w: save user: %w", ErrPersistence, err)
		}
	}

	return user, nil
}
