package noteshandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gonotes/internal/api"
	"gonotes/internal/lib/logger/slogdiscard"
	"gonotes/internal/middleware"
	"gonotes/internal/notes/dto"
	"gonotes/internal/notes/entity"
	"gonotes/internal/notes/noteshandler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

type mockNotesRepo struct {
	createFunc func(note *entity.Note) (*entity.Note, error)
	getFunc    func(id int, userID int) (*entity.Note, error)
	deleteFunc func(id int, userID int) (*entity.Note, error)
	getAllFunc func(userID int) ([]entity.Note, error)
}

func (m *mockNotesRepo) Create(n *entity.Note) (*entity.Note, error) {
	return m.createFunc(n)
}
func (m *mockNotesRepo) Get(id int, userID int) (*entity.Note, error) {
	return m.getFunc(id, userID)
}
func (m *mockNotesRepo) Delete(id, userID int) (*entity.Note, error) {
	return m.deleteFunc(id, userID)
}
func (m *mockNotesRepo) GetAll(userID int) ([]entity.Note, error) {
	return m.getAllFunc(userID)
}

func Test_Create(t *testing.T) {
	cases := []struct {
		name           string
		title          string
		content        string
		mockError      error
		expectedStatus int
		expectedError  string
		setupContext   bool
	}{
		{
			name:           "note created",
			title:          "test",
			content:        "test test test",
			expectedStatus: http.StatusCreated,
			setupContext:   true,
		},
		{
			name:           "invalid request body",
			title:          "",
			content:        "",
			expectedStatus: http.StatusBadRequest,
			setupContext:   true,
		},
		{
			name:           "repo error",
			title:          "fail",
			content:        "db fails",
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "db error",
			setupContext:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockNotesRepo{
				createFunc: func(n *entity.Note) (*entity.Note, error) {
					if tc.mockError != nil {
						return n, tc.mockError
					}
					return n, nil
				},
			}

			log := slogdiscard.NewDiscardLogger()
			h := noteshandler.NewHandler(log, repo)

			body, _ := json.Marshal(dto.NoteRequest{
				Title:   tc.title,
				Content: tc.content,
			})

			req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			if tc.setupContext {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()

			h.Create(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			if tc.expectedError != "" {
				var resp api.ErrResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.Equal(t, tc.expectedError, resp.Err)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	tests := []struct {
		name           string
		setupCtx       bool
		urlID          string
		mockGet        func(id int, userID int) (*entity.Note, error)
		expectedStatus int
	}{
		{
			name:           "unauthorized",
			setupCtx:       false,
			urlID:          "1",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid id format",
			setupCtx:       true,
			urlID:          "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "note not found",
			setupCtx: true,
			urlID:    "1",
			mockGet: func(id int, userID int) (*entity.Note, error) {
				return nil, errors.New("not found")
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "note retrieved",
			setupCtx: true,
			urlID:    "1",
			mockGet: func(id int, userID int) (*entity.Note, error) {
				return &entity.Note{ID: 1, Title: "test", Content: "content"}, nil
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockNotesRepo{getFunc: tc.mockGet}
			log := slogdiscard.NewDiscardLogger()
			h := noteshandler.NewHandler(log, repo)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.urlID)

			req := httptest.NewRequest(http.MethodGet, "/notes/"+tc.urlID, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tc.setupCtx {
				req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
			}

			rec := httptest.NewRecorder()
			h.Get(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				var resp dto.NoteResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.Equal(t, 1, resp.ID)
				assert.Equal(t, "test", resp.Title)
			}
			if res.StatusCode >= 400 {
				var resp api.ErrResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.NotEmpty(t, resp.Err)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	tests := []struct {
		name           string
		setupCtx       bool
		urlID          string
		mockDelete     func(id, userID int) (*entity.Note, error)
		expectedStatus int
	}{
		{
			name:           "unauthorized",
			setupCtx:       false,
			urlID:          "1",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid id format",
			setupCtx:       true,
			urlID:          "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "note not found",
			setupCtx: true,
			urlID:    "1",
			mockDelete: func(id, userID int) (*entity.Note, error) {
				return nil, errors.New("not found")
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "note deleted",
			setupCtx: true,
			urlID:    "1",
			mockDelete: func(id, userID int) (*entity.Note, error) {
				return &entity.Note{ID: 1, Title: "test", Content: "deleted"}, nil
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockNotesRepo{deleteFunc: tc.mockDelete}
			log := slogdiscard.NewDiscardLogger()
			h := noteshandler.NewHandler(log, repo)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.urlID)

			req := httptest.NewRequest(http.MethodDelete, "/notes/"+tc.urlID, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tc.setupCtx {
				req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
			}

			rec := httptest.NewRecorder()
			h.Delete(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				var resp dto.NoteResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.Equal(t, 1, resp.ID)
				assert.Equal(t, "test", resp.Title)
			}
			if res.StatusCode >= 400 {
				var resp api.ErrResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.NotEmpty(t, resp.Err)
			}
		})
	}
}

func Test_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		setupCtx       bool
		mockGetAll     func(userID int) ([]entity.Note, error)
		expectedStatus int
		expectedLen    int
	}{
		{
			name:           "unauthorized",
			setupCtx:       false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:     "repo error",
			setupCtx: true,
			mockGetAll: func(userID int) ([]entity.Note, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "notes retrieved",
			setupCtx: true,
			mockGetAll: func(userID int) ([]entity.Note, error) {
				return []entity.Note{
					{ID: 1, Title: "note1", Content: "content1"},
					{ID: 2, Title: "note2", Content: "content2"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedLen:    2,
		},
		{
			name:     "empty list",
			setupCtx: true,
			mockGetAll: func(userID int) ([]entity.Note, error) {
				return []entity.Note{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedLen:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockNotesRepo{getAllFunc: tc.mockGetAll}
			log := slogdiscard.NewDiscardLogger()
			h := noteshandler.NewHandler(log, repo)

			req := httptest.NewRequest(http.MethodGet, "/notes", nil)
			if tc.setupCtx {
				req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
			}
			rec := httptest.NewRecorder()

			h.GetAll(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				var resp []dto.NoteResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.Len(t, resp, tc.expectedLen)
			}
			if res.StatusCode >= 400 {
				var resp api.ErrResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.NotEmpty(t, resp.Err)
			}
		})
	}
}
