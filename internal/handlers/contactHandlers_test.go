package handlers

import (
    "io"
    "context"
    "bytes"
    "errors"
    "time"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/middleware"
)

type mockContactService struct {
	SaveContactFunc   func(ctx context.Context, contact internal.Contact) (int64, error)
	GetContactFunc    func(ctx context.Context, userID, contactID int64) (internal.Contact, error)
	ListContactsFunc  func(ctx context.Context, userID int64) ([]internal.Contact, error)
	UpdateContactFunc func(ctx context.Context, contact internal.Contact) (bool, error)
	DeleteContactFunc func(ctx context.Context, userID, contactID int64) (bool, error)
}


func (s *mockContactService) SaveContact(ctx context.Context, contact internal.Contact) (int64, error) {
	if s.SaveContactFunc != nil {
		return s.SaveContactFunc(ctx, contact)
	}
	return internal.ZERO, ErrSaveContactFuncNotImplemented
}

func (s *mockContactService) GetContact(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
	if s.GetContactFunc != nil {
		return s.GetContactFunc(ctx, userID, contactID)
	}
	return internal.Contact{}, ErrGetContactFuncNotImplemented
}

func (s *mockContactService) ListContacts(ctx context.Context, userID int64) ([]internal.Contact, error) {
	if s.ListContactsFunc != nil {
		return s.ListContactsFunc(ctx, userID)
	}
	return []internal.Contact{}, ErrListContactsFuncNotImplemented
}

func (s *mockContactService) UpdateContact(ctx context.Context, contact internal.Contact) (bool, error) {
	if s.UpdateContactFunc != nil {
		return s.UpdateContactFunc(ctx, contact)
	}
	return false, ErrUpdateContactFuncNotImplemented
}

func (s *mockContactService) DeleteContact(ctx context.Context, userID, contactID int64) (bool, error) {
	if s.DeleteContactFunc != nil {
		return s.DeleteContactFunc(ctx, userID, contactID)
	}
}

func TestContactHandler_SaveContact(t *testing.T) {
    mockService := &mockContactService {
        SaveContactFunc: func(ctx context.Context, contact internal.Contact) (int64, error) {
            if contact.UserID <= internal.ZERO {
                return internal.ZERO, ErrMissingUserID
            }
            if contact.Name == internal.EMPTY {
                return internal.ZERO, ErrMissingContactName
            }
            if contact.DDI == internal.EMPTY {
                return internal.ZERO, ErrMissingContactDDI
            }
            if contact.DDD == internal.EMPTY {
                return internal.ZERO, ErrMissingContactDDD
            }
            if contact.Name == internal.EMPTY {
                return internal.ZERO, ErrMissingContactPhoneNumber
            }

            return 1, nil
        },
    }

    handler := NewContactHandler(mockService)

    tests := []struct {
        name            string
        userID          int64
        inputBody       string
        wantStatus      int
        wantResp        string
   }{
       {
           name:        "Contato salvo com sucesso",
           userID:      1,
           inputBody:   `{"id": ', "userid": 1, "name:" "john", "ddi": "055", "ddd": "051", "phonenumber": "123456789"}`,
           wantStatus:  http.StatusCreated,
           wantResp:    `{"response": 1}`,
       }
   }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, "/save-contact", bytes.NewBufferString(tt.inputBody))
            req.Header.Set("Content-Type", "application/json")
            ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
            req = req.WithContext(ctx)
            rec := httptest.NewRecorder()
            handler.SaveContact(rec, req)
            res := rec.Result()
            defer res.Body.Close()
            if res.StatusCode != tt.wantStatus {
                t.Errorf("[%s] Status esperado: %d, recebido: %d", tt.name, tt.wantStatus, res.StatusCode)
            }

            var respBody bytes.Buffer
            respBody.ReadFrom(res.Body)
            respStr := respBody.String()
            if respStr != tt.wantResp+"\n" {
                t.Errorf("[%s] Resposta esperada: %s, recebida: %s", tt.name, tt.wantResp, respStr)
           }
       })
   }
}
