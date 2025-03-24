package handlers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amarantec/move-easy/internal"
	"github.com/amarantec/move-easy/internal/middleware"
)

// Mock do IAddressService
type mockAddressService struct {
	GetAddressFunc         func(ctx context.Context, userID int64) (internal.Address, error)
	AddOrUpdateAddressFunc func(ctx context.Context, address internal.Address) (int64, error)
}

func (m *mockAddressService) GetAddress(ctx context.Context, userID int64) (internal.Address, error) {
	return m.GetAddressFunc(ctx, userID)
}

func (m *mockAddressService) AddOrUpdateAddress(ctx context.Context, address internal.Address) (int64, error) {
	return m.AddOrUpdateAddressFunc(ctx, address)
}

// Test do handler GetAddress
func TestAddressHandler_GetAddress(t *testing.T) {
	mockService := &mockAddressService{
		GetAddressFunc: func(ctx context.Context, userID int64) (internal.Address, error) {
			if userID <= internal.ZERO {
				return internal.Address{}, ErrMissingUserID
			}
			addr := internal.Address{
				ID:           1,
				UserID:       1,
				Street:       "General Osório",
				Number:       "2211",
				CEP:          "95520000",
				Neighborhood: "Glória",
				City:         "Osório",
				State:        "RS",
				CreatedAt:    time.Time{},
				UpdatedAt:    nil,
				DeletedAt:    nil,
			}
			return addr, nil
		},
	}

	handler := NewAddressHandler(mockService)

	tests := []struct {
		name       string
		userID     int64
		wantStatus int
		wantResp   string
	}{
		{
			name:       "Address Get Successfully",
			userID:     1,
			wantStatus: http.StatusOK,
			wantResp:   `{"response":{"ID":1,"UserID":1,"Street":"General Osório","Number":"2211","CEP":"95520000","Neighborhood":"Glória","City":"Osório","State":"RS","CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":null,"DeletedAt":null}}`,
		},
		{
			name:       "Missing userID",
			userID:     0,
			wantStatus: http.StatusInternalServerError,
			wantResp:   "could not get this address, error: missing user id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/get-address", nil)

			ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)

			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetAddress(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("[%s] Esperado status %d, recebeu %d", tt.name, tt.wantStatus, res.StatusCode)
			}

			body, _ := io.ReadAll(res.Body)
			respStr := string(body)

			if respStr != tt.wantResp+"\n" {
				t.Errorf("[%s] Resposta esperada: %s, recebeu: %s", tt.name, tt.wantResp, respStr)
			}
		})
	}
}

// Teste do handler AddOrUpdateAddress
func TestAddressHandler_AddOrUpdateAddress(t *testing.T) {
	mockService := &mockAddressService{
		AddOrUpdateAddressFunc: func(ctx context.Context, address internal.Address) (int64, error) {
			if address.UserID <= internal.ZERO {
				return internal.ZERO, ErrMissingUserID
			}

			if address.Street == internal.EMPTY {
				return internal.ZERO, ErrMissingAddressStreet
			}

			if address.Number == internal.EMPTY {
				return internal.ZERO, ErrMissingAddressNumber
			}

			if address.CEP == internal.EMPTY {
				return internal.ZERO, ErrMissingAddressCEP
			}

			if address.Neighborhood == internal.EMPTY {
				return internal.ZERO, ErrMissingAddressNeighborhood
			}

			if address.City == internal.EMPTY {
				return internal.ZERO, ErrMissingAddressCity
			}

			if address.State == internal.EMPTY {
				return internal.ZERO, ErrMissingAddressState
			}

			return 1, nil
		},
	}

	handler := NewAddressHandler(mockService)

	tests := []struct {
		name       string
		userID     int64
		inputBody  string
		wantStatus int
		wantResp   string
	}{
		{
			name:       "Successful insert",
			userID:     1,
			inputBody:  `{"id": 1, "userid": 1, "street": "General Osório", "number": "2211", "cep": "95520000", "neighborhood": "Glória", "city": "Osório", "state": "RS"}`,
			wantStatus: http.StatusCreated,
			wantResp:   `{"response":1}`,
		},
		{
			name:       "Missing User ID",
			userID:     0,
			inputBody:  `{"id": 1, "userid": 0, "street": "General Osório", "number": "2211", "cep": "95520000", "neighborhood": "Glória", "city": "Osório", "state": "RS"}`,
			wantStatus: http.StatusInternalServerError,
			wantResp:   "could not save this address, error: missing user id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/save-address", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			handler.AddOrUpdateAddress(rec, req)
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

var ErrMissingUserID = errors.New("missing user id")
var ErrMissingAddressStreet = errors.New("missing address street")
var ErrMissingAddressNumber = errors.New("missing address number")
var ErrMissingAddressCEP = errors.New("missing address cep")
var ErrMissingAddressNeighborhood = errors.New("missing address neighborhood")
var ErrMissingAddressCity = errors.New("missing address city")
var ErrMissingAddressState = errors.New("missing address state")
