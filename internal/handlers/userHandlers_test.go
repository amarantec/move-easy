package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amarantec/move-easy/internal"
)

// Mock do IUserService
type mockUserService struct {
	RegisterFunc            func(ctx context.Context, user internal.UserRegister) (int64, error)
	ValidateCredentialsFunc func(ctx context.Context, user internal.UserLogin) (string, error)
}

func (m *mockUserService) Register(ctx context.Context, user internal.UserRegister) (int64, error) {
	return m.RegisterFunc(ctx, user)
}

func (m *mockUserService) ValidateCredentials(ctx context.Context, user internal.UserLogin) (string, error) {
	return m.ValidateCredentialsFunc(ctx, user)
}

// Teste do handler Register
func TestUserHandler_Register(t *testing.T) {
	mockService := &mockUserService{
		RegisterFunc: func(ctx context.Context, user internal.UserRegister) (int64, error) {
			if user.Email == internal.EMPTY {
                return internal.ZERO, ErrMissingEmail
            }
            if user.Password == internal.EMPTY {
				return internal.ZERO, ErrMissingPassword
			}
			return 1, nil
		},
	}

	handler := NewUserHandler(mockService)

	tests := []struct {
		name       string
		inputBody  string
		wantStatus int
		wantResp   string
	}{
		{
			name:       "Successful registration",
			inputBody:  `{"email": "john@example.com", "password": "securepass"}`,
			wantStatus: http.StatusCreated,
			wantResp:   `{"response":1}`,
		},
		{
			name:       "Missing email",
			inputBody:  `{"email": "", "password": "securepass"}`,
			wantStatus: http.StatusInternalServerError,
			wantResp:   `could not register this user, error: email is required`,
		},
		{
			name:       "Missing password",
			inputBody:  `{"email": "john@example.com", "password": ""}`,
			wantStatus: http.StatusInternalServerError,
			wantResp:   `could not register this user, error: password is required`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Register(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("[%s] Esperado status %d, recebeu %d", tt.name, tt.wantStatus, res.StatusCode)
			}

			var respBody bytes.Buffer
			respBody.ReadFrom(res.Body)
			respStr := respBody.String()

			if respStr != tt.wantResp+"\n" {
				t.Errorf("[%s] Resposta esperada: %s, recebeu: %s", tt.name, tt.wantResp, respStr)
			}
		})
	}
}

// Teste do handler Login
func TestUserHandler_Login(t *testing.T) {
	mockService := &mockUserService{
		ValidateCredentialsFunc: func(ctx context.Context, user internal.UserLogin) (string, error) {
			if user.Email == "invalid@example.com" {
				return internal.EMPTY, ErrInvalidCredentials
			}
			return "mocked_token_123", nil
		},
	}

	handler := NewUserHandler(mockService)

	tests := []struct {
		name       string
		inputBody  string
		wantStatus int
		wantResp   string
	}{
		{
			name:       "Successful login",
			inputBody:  `{"email": "john@example.com", "password": "securepass"}`,
			wantStatus: http.StatusOK,
			wantResp:   `{"token":"mocked_token_123"}`,
		},
		{
			name:       "Invalid credentials",
			inputBody:  `{"email": "invalid@example.com", "password": "wrongpass"}`,
			wantStatus: http.StatusInternalServerError,
			wantResp:   `could not validate this credentials, error: invalid credentials`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Login(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("[%s] Esperado status %d, recebeu %d", tt.name, tt.wantStatus, res.StatusCode)
			}

			var respBody bytes.Buffer
			respBody.ReadFrom(res.Body)
			respStr := respBody.String()

			if respStr != tt.wantResp+"\n" {
				t.Errorf("[%s] Resposta esperada: %s, recebeu: %s", tt.name, tt.wantResp, respStr)
			}
		})
	}
}

var ErrMissingEmail = errors.New("email is required")
var ErrMissingPassword = errors.New("password is required")
var ErrInvalidCredentials = errors.New("invalid credentials")
