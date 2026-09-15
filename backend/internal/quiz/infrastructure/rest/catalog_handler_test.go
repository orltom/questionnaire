package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"uuid"

	"go.uber.org/mock/gomock"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/domain"
)

func TestCatalogHandler_CreateCatalog(t *testing.T) {
	defaultCatalog, _ := domain.NewCatalog("Animals", "Questions about animals")

	type fields struct {
		service func(s *MockcatalogService)
	}
	type args struct {
		body string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       *api.Catalog
	}{
		{
			name: "When the request is valid, then create the catalog and return it",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Create(gomock.Any(), "Animals", "Questions about animals").Return(defaultCatalog, nil)
				},
			},
			args: args{
				body: `{"title":"Animals","description":"Questions about animals"}`,
			},
			wantStatus: http.StatusCreated,
			want: &api.Catalog{
				Id:          api.CatalogId(defaultCatalog.ID()),
				Title:       "Animals",
				Description: "Questions about animals",
				QuestionIds: []uuid.UUID{},
			},
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockcatalogService) {},
			},
			args: args{
				body: `{"title":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the catalog is invalid, then return bad request",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Catalog{}, application.ErrInvalidArguments)
				},
			},
			args: args{
				body: `{"title":"","description":""}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Catalog{}, errors.New("database is down"))
				},
			},
			args: args{
				body: `{"title":"Animals","description":"Questions about animals"}`,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockcatalogService(ctrl)
			tt.fields.service(service)

			h := catalogHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/catalogs", strings.NewReader(tt.args.body))

			h.CreateCatalog(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("CreateCatalog() status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.want == nil {
				return
			}

			var got api.Catalog
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("CreateCatalog() response body is not a catalog: %v", err)
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("CreateCatalog() = %+v, want %+v", got, *tt.want)
			}
		})
	}
}

func TestCatalogHandler_GetCatalog(t *testing.T) {
	defaultCatalog, _ := domain.NewCatalog("Animals", "Questions about animals")

	type fields struct {
		service func(s *MockcatalogService)
	}
	type args struct {
		catalogID api.CatalogId
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       *api.Catalog
	}{
		{
			name: "When the catalog exists, then return it",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Get(gomock.Any(), defaultCatalog.ID()).Return(defaultCatalog, nil)
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
			},
			wantStatus: http.StatusOK,
			want: &api.Catalog{
				Id:          api.CatalogId(defaultCatalog.ID()),
				Title:       "Animals",
				Description: "Questions about animals",
				QuestionIds: []uuid.UUID{},
			},
		},
		{
			name: "When the catalog does not exist, then return not found",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.Catalog{}, application.ErrEntityNotFound)
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.Catalog{}, errors.New("database is down"))
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockcatalogService(ctrl)
			tt.fields.service(service)

			h := catalogHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/catalogs/"+tt.args.catalogID.String(), nil)

			h.GetCatalog(rec, req, tt.args.catalogID)

			if rec.Code != tt.wantStatus {
				t.Fatalf("GetCatalog() status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.want == nil {
				return
			}

			var got api.Catalog
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("GetCatalog() response body is not a catalog: %v", err)
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("GetCatalog() = %+v, want %+v", got, *tt.want)
			}
		})
	}
}

func TestCatalogHandler_UpdateCatalog(t *testing.T) {
	defaultCatalog, _ := domain.NewCatalog("Animals", "Questions about animals")

	type fields struct {
		service func(s *MockcatalogService)
	}
	type args struct {
		catalogID api.CatalogId
		body      string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "When the request is valid, then update the catalog",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Update(gomock.Any(), defaultCatalog.ID(), "Plants", "Questions about plants").Return(nil)
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
				body:      `{"title":"Plants","description":"Questions about plants"}`,
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the request body is malformed, then return bad request",
			fields: fields{
				service: func(s *MockcatalogService) {},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
				body:      `{"title":`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the catalog is invalid, then return bad request",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrInvalidArguments)
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
				body:      `{"title":"","description":""}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "When the catalog does not exist, then return not found",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
				body:      `{"title":"Plants","description":"Questions about plants"}`,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
				body:      `{"title":"Plants","description":"Questions about plants"}`,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockcatalogService(ctrl)
			tt.fields.service(service)

			h := catalogHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/catalogs/"+tt.args.catalogID.String(), strings.NewReader(tt.args.body))

			h.UpdateCatalog(rec, req, tt.args.catalogID)

			if rec.Code != tt.wantStatus {
				t.Errorf("UpdateCatalog() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestCatalogHandler_DeleteCatalog(t *testing.T) {
	defaultCatalog, _ := domain.NewCatalog("Animals", "Questions about animals")

	type fields struct {
		service func(s *MockcatalogService)
	}
	type args struct {
		catalogID api.CatalogId
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "When the catalog exists, then delete it",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Delete(gomock.Any(), defaultCatalog.ID()).Return(nil)
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "When the catalog does not exist, then return not found",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(application.ErrEntityNotFound)
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "When the service fails, then return internal server error",
			fields: fields{
				service: func(s *MockcatalogService) {
					s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("database is down"))
				},
			},
			args: args{
				catalogID: api.CatalogId(defaultCatalog.ID()),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockcatalogService(ctrl)
			tt.fields.service(service)

			h := catalogHandler{service: service}
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/catalogs/"+tt.args.catalogID.String(), nil)

			h.DeleteCatalog(rec, req, tt.args.catalogID)

			if rec.Code != tt.wantStatus {
				t.Errorf("DeleteCatalog() status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
