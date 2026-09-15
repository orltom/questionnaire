package identity

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestUserService_Authenticate(t *testing.T) {
	existing, _ := NewUser("https://accounts.google.com", "1234", "Ada")

	type mocks struct {
		repository func(r *MockUserRepository)
	}
	type args struct {
		issuer  string
		subject string
		name    string
	}
	tests := []struct {
		name     string
		mocks    mocks
		args     args
		wantName string
		wantErr  error
	}{
		{
			name: "When the user is known, then return the existing user",
			mocks: mocks{
				repository: func(r *MockUserRepository) {
					r.EXPECT().FindBySubject(gomock.Any(), "https://accounts.google.com", "1234").Return(existing, nil)
				},
			},
			args: args{
				issuer:  "https://accounts.google.com",
				subject: "1234",
				name:    "Ada",
			},
			wantName: "Ada",
			wantErr:  nil,
		},
		{
			name: "When the user is unknown, then create it",
			mocks: mocks{
				repository: func(r *MockUserRepository) {
					r.EXPECT().FindBySubject(gomock.Any(), gomock.Any(), gomock.Any()).Return(User{}, ErrEntityNotFound)
					r.EXPECT().Create(gomock.Any(), gomock.Cond(func(u User) bool {
						return u.Issuer() == "https://accounts.google.com" && u.Subject() == "5678"
					})).Return(nil)
				},
			},
			args: args{
				issuer:  "https://accounts.google.com",
				subject: "5678",
				name:    "Grace",
			},
			wantName: "Grace",
			wantErr:  nil,
		},
		{
			name: "When the issuer is missing, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockUserRepository) {
					r.EXPECT().FindBySubject(gomock.Any(), gomock.Any(), gomock.Any()).Return(User{}, ErrEntityNotFound)
				},
			},
			args: args{
				issuer:  "",
				subject: "5678",
				name:    "Grace",
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockUserRepository) {
					r.EXPECT().FindBySubject(gomock.Any(), gomock.Any(), gomock.Any()).Return(User{}, errors.New("database is down"))
				},
			},
			args: args{
				issuer:  "https://accounts.google.com",
				subject: "5678",
				name:    "Grace",
			},
			wantErr: ErrPersistence,
		},
		{
			name: "When saving the user fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockUserRepository) {
					r.EXPECT().FindBySubject(gomock.Any(), gomock.Any(), gomock.Any()).Return(User{}, ErrEntityNotFound)
					r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				issuer:  "https://accounts.google.com",
				subject: "5678",
				name:    "Grace",
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockUserRepository(ctrl)
			tt.mocks.repository(repo)

			s := NewUserService(repo)
			got, err := s.Authenticate(context.Background(), tt.args.issuer, tt.args.subject, tt.args.name)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.Name() != tt.wantName {
				t.Errorf("Authenticate() name = %q, want %q", got.Name(), tt.wantName)
			}
		})
	}
}
