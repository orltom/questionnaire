package identity

import (
	"context"
	"errors"
	"fmt"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindBySubject(ctx context.Context, issuer string, subject string) (User, error)
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
	user, err := s.repository.FindBySubject(ctx, issuer, subject)
	if err == nil {
		if user.Name() != name && len(name) > 0 {
			user.Rename(name)
		}

		return user, nil
	}

	if !errors.Is(err, ErrEntityNotFound) {
		return User{}, fmt.Errorf("%w: load user: %w", ErrPersistence, err)
	}

	user, err = NewUser(issuer, subject, name)
	if err != nil {
		return User{}, fmt.Errorf("%w: create user: %w", ErrInvalidArguments, err)
	}

	err = s.repository.Create(ctx, user)
	if err != nil {
		return User{}, fmt.Errorf("%w: save user: %w", ErrPersistence, err)
	}

	return user, nil
}
