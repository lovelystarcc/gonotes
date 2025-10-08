package authhandler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gonotes/internal/api"
	"gonotes/internal/auth"
	"gonotes/internal/auth/authhandler"
	"gonotes/internal/auth/dto"
	"gonotes/internal/auth/entity"
	"gonotes/internal/lib/logger/slogdiscard"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	createFunc func(u *entity.User) (int, error)
	getFunc    func(email string) (*entity.User, error)
	checkFunc  func(id, something int) bool
	deleteFunc func(id int) error
}

func (m *mockUserRepo) Create(u *entity.User) (int, error) {
	return m.createFunc(u)
}
func (m *mockUserRepo) Get(e string) (*entity.User, error) {
	return m.getFunc(e)
}
func (m *mockUserRepo) Check(id, smth int) bool {
	return m.checkFunc(id, smth)
}
func (m *mockUserRepo) Delete(id int) error {
	return m.deleteFunc(id)
}

func Test_Create(t *testing.T) {
	cases := []struct {
		name           string
		email          string
		password       string
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "user created",
			email:          "test@example.com",
			password:       "12345",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "user already exists",
			email:          "exists@example.com",
			password:       "12345",
			mockError:      auth.ErrUserExists,
			expectedStatus: http.StatusConflict,
			expectedError:  auth.ErrUserExists.Error(),
		},
		{
			name:           "invalid request body",
			email:          "",
			password:       "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockUserRepo{
				createFunc: func(u *entity.User) (int, error) {
					if tc.mockError != nil {
						return 0, tc.mockError
					}
					return 1, nil
				},
			}

			log := slogdiscard.NewDiscardLogger()
			h := authhandler.NewHandler(log, repo, "secret")

			body, _ := json.Marshal(dto.UserRequest{
				Email:    tc.email,
				Password: tc.password,
			})

			req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.Create(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			var resp api.ErrResponse
			_ = json.NewDecoder(res.Body).Decode(&resp)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, resp.Err)
			}
		})
	}
}

func Test_Login(t *testing.T) {
	cases := []struct {
		name           string
		email          string
		password       string
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "user logged in",
			email:          "test@example.com",
			password:       "12345",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			email:          "",
			password:       "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "email is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockUserRepo{
				getFunc: func(email string) (*entity.User, error) {
					hashed, _ := bcrypt.GenerateFromPassword([]byte(tc.password), bcrypt.DefaultCost)
					return &entity.User{ID: 1, Email: email, Password: string(hashed)}, tc.mockError
				},
			}

			log := slogdiscard.NewDiscardLogger()
			h := authhandler.NewHandler(log, repo, "secret")

			body, _ := json.Marshal(dto.UserRequest{
				Email:    tc.email,
				Password: tc.password,
			})

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Login(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if tc.expectedStatus == http.StatusOK {
				var resp dto.LoginResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.Equal(t, tc.email, resp.Email)
				assert.NotEmpty(t, resp.Token)
			} else {
				var resp api.ErrResponse
				_ = json.NewDecoder(res.Body).Decode(&resp)
				assert.Equal(t, tc.expectedError, resp.Err)
			}

		})
	}
}
