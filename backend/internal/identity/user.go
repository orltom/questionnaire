package identity

import (
	"errors"
	"uuid"
)

type UserID uuid.UUID

type User struct {
	id      UserID
	issuer  string
	subject string
	name    string
}

func NewUser(issuer string, subject string, name string) (User, error) {
	if len(issuer) == 0 {
		return User{}, errors.New("issuer is required")
	}

	if len(subject) == 0 {
		return User{}, errors.New("subject is required")
	}

	return User{
		id:      UserID(uuid.NewV7()),
		issuer:  issuer,
		subject: subject,
		name:    name,
	}, nil
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Issuer() string {
	return u.issuer
}

func (u *User) Subject() string {
	return u.subject
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Rename(name string) {
	u.name = name
}
