package application

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestCatalogService_Create(t *testing.T) {
	defaultCatalog, _ := domain.NewCatalog("Animals", "Questions about animals")

	type mocks struct {
		repository func(r *MockCatalogRepository)
	}
	type args struct {
		title       string
		description string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    domain.Catalog
		wantErr error
	}{
		{
			name: "When creating a valid catalog, then save it",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Cond(func(c domain.Catalog) bool {
						return c.Title() == "Animals" &&
							c.Description() == "Questions about animals"
					})).Return(nil)
				},
			},

			args: args{
				title:       "Animals",
				description: "Questions about animals",
			},
			want:    defaultCatalog,
			wantErr: nil,
		},
		{
			name: "When the title is empty, then return an invalid arguments error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {},
			},
			args: args{
				title:       "",
				description: "",
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				title:       "Animals",
				description: "Questions about animals",
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockCatalogRepository(ctrl)
			tt.mocks.repository(repo)

			s := &CatalogService{repository: repo, lookupService: nil}
			got, err := s.Create(context.Background(), tt.args.title, tt.args.description)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(domain.Catalog{}), cmpopts.IgnoreFields(domain.Catalog{}, "id")); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCatalogService_Get(t *testing.T) {
	defaultCatalog, _ := domain.NewCatalog("Animals", "Questions about animals")

	type mocks struct {
		repository func(r *MockCatalogRepository)
	}
	type args struct {
		id domain.CatalogID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    domain.Catalog
		wantErr error
	}{
		{
			name: "When the catalog exists, then return it",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Find(gomock.Any(), defaultCatalog.ID()).Return(defaultCatalog, nil)
				},
			},

			args: args{
				id: defaultCatalog.ID(),
			},
			want:    defaultCatalog,
			wantErr: nil,
		},
		{
			name: "When the catalog ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Find(gomock.Any(), defaultCatalog.ID()).Return(domain.Catalog{}, ErrEntityNotFound)
				},
			},
			args: args{
				id: defaultCatalog.ID(),
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Find(gomock.Any(), defaultCatalog.ID()).Return(domain.Catalog{}, errors.New("database is down"))
				},
			},
			args: args{
				id: defaultCatalog.ID(),
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockCatalogRepository(ctrl)
			tt.mocks.repository(repo)

			s := &CatalogService{repository: repo, lookupService: nil}
			got, err := s.Get(context.Background(), tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(domain.Catalog{})); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCatalogService_Delete(t *testing.T) {
	id := domain.CatalogID(uuid.Nil())

	type mocks struct {
		repository func(r *MockCatalogRepository)
	}
	type args struct {
		id domain.CatalogID
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When the catalog exists, then delete it",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Delete(gomock.Any(), id).Return(nil)
				},
			},

			args: args{
				id: id,
			},
			wantErr: nil,
		},
		{
			name: "When the catalog ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Delete(gomock.Any(), id).Return(ErrEntityNotFound)
				},
			},
			args: args{
				id: id,
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Delete(gomock.Any(), id).Return(errors.New("database is down"))
				},
			},
			args: args{
				id: id,
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockCatalogRepository(ctrl)
			tt.mocks.repository(repo)

			s := &CatalogService{repository: repo, lookupService: nil}
			err := s.Delete(context.Background(), tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCatalogService_Update(t *testing.T) {
	defaultCatalog, _ := domain.NewCatalog("Sport", "Questions about animals")

	type mocks struct {
		repository func(r *MockCatalogRepository)
	}
	type args struct {
		id          domain.CatalogID
		title       string
		description string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr error
	}{
		{
			name: "When updating an existing catalog, then save the changes",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Find(gomock.Any(), defaultCatalog.ID()).Return(defaultCatalog, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Cond(func(c domain.Catalog) bool {
						return c.Title() == "Sport" && c.Description() == "Chess is a sport"
					})).Return(nil)
				},
			},

			args: args{
				id:          defaultCatalog.ID(),
				title:       "Sport",
				description: "Chess is a sport",
			},
			wantErr: nil,
		},
		{
			name: "When the catalog ID is unknown, then return an entity not found error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Find(gomock.Any(), defaultCatalog.ID()).Return(domain.Catalog{}, ErrEntityNotFound)
				},
			},
			args: args{
				id:          defaultCatalog.ID(),
				title:       "Sport",
				description: "Chess is a sport",
			},
			wantErr: ErrEntityNotFound,
		},
		{
			name: "When updated fields are not valid, then return an invalid argument error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Find(gomock.Any(), defaultCatalog.ID()).Return(defaultCatalog, nil)
				},
			},
			args: args{
				id:          defaultCatalog.ID(),
				title:       "",
				description: "",
			},
			wantErr: ErrInvalidArguments,
		},
		{
			name: "When the repository fails, then return a persistence error",
			mocks: mocks{
				repository: func(r *MockCatalogRepository) {
					r.EXPECT().Find(gomock.Any(), defaultCatalog.ID()).Return(defaultCatalog, nil)
					r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				id:          defaultCatalog.ID(),
				title:       "Sport",
				description: "Chess is a sport",
			},
			wantErr: ErrPersistence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockCatalogRepository(ctrl)
			tt.mocks.repository(repo)

			s := &CatalogService{repository: repo, lookupService: nil}
			err := s.Update(context.Background(), tt.args.id, tt.args.title, tt.args.description)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
