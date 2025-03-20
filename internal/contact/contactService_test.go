package contact

import (
	"context"
	"errors"
	"github.com/amarantec/move-easy/internal"
	"testing"
    "reflect"
)

type mockContactRepository struct {
	SaveContactFunc   func(ctx context.Context, contact internal.Contact) (int64, error)
	GetContactFunc    func(ctx context.Context, userID, contactID int64) (internal.Contact, error)
	ListContactsFunc  func(ctx context.Context, userID int64) ([]internal.Contact, error)
	UpdateContactFunc func(ctx context.Context, contact internal.Contact) (bool, error)
	DeleteContactFunc func(ctx context.Context, userID, contactID int64) (bool, error)
}

func (m *mockContactRepository) SaveContact(ctx context.Context, contact internal.Contact) (int64, error) {
	if m.SaveContactFunc != nil {
		return m.SaveContactFunc(ctx, contact)
	}
	return internal.ZERO, ErrSaveContactFuncNotImplemented
}

func (m *mockContactRepository) GetContact(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
	if m.GetContactFunc != nil {
		return m.GetContactFunc(ctx, userID, contactID)
	}
	return internal.Contact{}, ErrGetContactFuncNotImplemented
}

func (m *mockContactRepository) ListContacts(ctx context.Context, userID int64) ([]internal.Contact, error) {
	if m.ListContactsFunc != nil {
		return m.ListContactsFunc(ctx, userID)
	}
	return []internal.Contact{}, ErrListContactsFuncNotImplemented
}

func (m *mockContactRepository) UpdateContact(ctx context.Context, contact internal.Contact) (bool, error) {
	if m.UpdateContactFunc != nil {
		return m.UpdateContactFunc(ctx, contact)
	}
	return false, ErrUpdateContactFuncNotImplemented
}

func (m *mockContactRepository) DeleteContact(ctx context.Context, userID, contactID int64) (bool, error) {
	if m.DeleteContactFunc != nil {
		return m.DeleteContactFunc(ctx, userID, contactID)
	}
	return false, ErrDeleteContactFuncNotImplemented
}

func TestSaveContact(t *testing.T) {
	tests := []struct {
		name         string
		userID       int64
		input        internal.Contact
		mockFunc     func(ctx context.Context, contact internal.Contact) (int64, error)
		wantResponse int64
		wantError    bool
	}{
		{
			name:   "Contact saved successfully",
			userID: 1,
			input:  internal.Contact{UserID: 1, Name: "Test", DDI: "055", DDD: "051", PhoneNumber: "123456789"},
			mockFunc: func(ctx context.Context, contact internal.Contact) (int64, error) {
				return 1, nil
			},
			wantResponse: 1,
			wantError:    false,
		},
		{
			name:   "Err Missing Fields",
			userID: 1,
			input:  internal.Contact{UserID: 1, Name: "", DDI: "055", DDD: "051", PhoneNumber: "123456789"},
			mockFunc: func(ctx context.Context, contact internal.Contact) (int64, error) {
				return internal.ZERO, ErrInvalidContact
			},
			wantResponse: internal.ZERO,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockContactRepository{
				SaveContactFunc: tt.mockFunc,
			}
			service := NewContactService(mockRepo)
			id, err := service.SaveContact(context.Background(), tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v\n", tt.name, tt.wantError, err)
			}

			if id != tt.wantResponse {
				t.Errorf("[%s] ID esperado: %d, recebido: %d\n", tt.name, tt.wantResponse, id)
			}
		})
	}
}

func TestGetContact(t *testing.T) {
    tests := []struct {
        name        string
        userID      int64
        contactID   int64
        mockFunc    func(ctx context.Context, userID, contactID int64) (internal.Contact, error)
        wantResp    internal.Contact
        wantErr     bool
    }{
        {
            name:       "Contato encontrado",
            userID:     1,
            contactID:  1,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
                return internal.Contact{
                    ID:          1,
                    UserID:      1,
                    Name:        "John",
                    DDI:         "055",
                    DDD:         "051",
                    PhoneNumber: "123456789",
                }, nil
            },
            wantResp: internal.Contact{
                ID:          1,
                UserID:      1,
                Name:        "John",
                DDI:         "055",
                DDD:         "051",
                PhoneNumber: "123456789",
            },
            wantErr: false,
        },
        {
            name:       "Contato não encontrado",
            userID:     1,
            contactID:  2,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
                return internal.Contact{}, nil
            },
            wantResp: internal.Contact{},
            wantErr:  false,
       },
        {
            name:       "Contact ID vazio",
            userID:     1,
            contactID:  internal.ZERO,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
                return internal.Contact{}, ErrContactIDEmpty
            },
            wantResp: internal.Contact{},
            wantErr:  true,
       },
        {
            name:       "UserID vazio",
            userID:     internal.ZERO,
            contactID:  internal.ZERO,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
                return internal.Contact{}, ErrUserIDEmpty
            },
            wantResp: internal.Contact{},
            wantErr:  true,
       },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockContactRepository{
                GetContactFunc: tt.mockFunc,
            }
            service := NewContactService(mockRepo)

            contact, err := service.GetContact(context.Background(), tt.userID, tt.contactID)
            if (err != nil) != tt.wantErr {
                t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v", tt.name, tt.wantErr, err)
            }
            if contact != tt.wantResp {
                t.Errorf("[%s] Contato esperado: %+v, recebido: %+v", tt.name, tt.wantResp, contact)
            }
        })
    }
}

func TestListContacts(t *testing.T) {
    tests := []struct {
        name        string
        userID      int64
        mockFunc    func(ctx context.Context, userID int64) ([]internal.Contact, error)
        wantResp    []internal.Contact
        wantErr     bool
    }{
        {
            name:       "Lista de contatos retornada com sucesso",
            userID:     1,
            mockFunc:   func(ctx context.Context, userID int64) ([]internal.Contact, error) {
                return []internal.Contact{
                    {  
                        ID:             1,
                        UserID:         1,
                        Name:           "John",
                        DDI:            "055",
                        DDD:            "051",
                        PhoneNumber:    "123456789",
                    },
                }, nil
            },
            wantResp: []internal.Contact{
                {
                    ID:          1,
                    UserID:      1,
                    Name:        "John",
                    DDI:         "055",
                    DDD:         "051",
                    PhoneNumber: "123456789",
                },
            },
            wantErr: false,
        },
        {
            name:       "UserID vazio",
            userID:     internal.ZERO,
            mockFunc:   func(ctx context.Context, userID int64) ([]internal.Contact, error) {
                return []internal.Contact{}, ErrUserIDEmpty
            },
            wantResp: []internal.Contact{},
            wantErr:  true,
       },
       {
           name:        "Lista de contatos vazia",
           userID:      1,
           mockFunc:    func(ctx context.Context, userID int64) ([]internal.Contact, error) {
               return []internal.Contact{}, nil
           },
           wantResp:    []internal.Contact{},
           wantErr:     false,
       },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockContactRepository{
                ListContactsFunc: tt.mockFunc,
            }
            service := NewContactService(mockRepo)

            contacts, err := service.ListContacts(context.Background(), tt.userID)
            if (err != nil) != tt.wantErr {
                t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v", tt.name, tt.wantErr, err)
            }
            if !reflect.DeepEqual(contacts, tt.wantResp) {
                t.Errorf("[%s] Contato esperado: %+v, recebido: %+v", tt.name, tt.wantResp, contacts)
            }
        })
    }
}

func TestUpdateContact(t *testing.T) {
	tests := []struct {
		name         string
		input        internal.Contact
		mockFunc     func(ctx context.Context, contact internal.Contact) (bool, error)
		wantResponse bool
		wantError    bool
	}{
		{
			name:   "Contato atualizado",
			input:  internal.Contact{ID: 1, UserID: 1, Name: "Test", DDI: "055", DDD: "051", PhoneNumber: "123456789"},
			mockFunc: func(ctx context.Context, contact internal.Contact) (bool, error) {
				return true, nil
			},
			wantResponse: true,
			wantError:    false,
		},
		{
			name:   "Erro - faltando parâmetros",
			input:  internal.Contact{ID: 1, UserID: 1, Name: "", DDI: "055", DDD: "051", PhoneNumber: "123456789"},
			mockFunc: func(ctx context.Context, contact internal.Contact) (bool, error) {
				return false, ErrInvalidContact
			},
			wantResponse: false,
			wantError:    true,
		},
		{
			name:   "Contato não encontrado",
			input:  internal.Contact{ID: 2, UserID: 1, Name: "Test", DDI: "055", DDD: "051", PhoneNumber: "123456789"},
			mockFunc: func(ctx context.Context, contact internal.Contact) (bool, error) {
				return false, ErrContactNotFound
			},
			wantResponse: false,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockContactRepository{
				UpdateContactFunc: tt.mockFunc,
			}
			service := NewContactService(mockRepo)
			response, err := service.UpdateContact(context.Background(), tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v\n", tt.name, tt.wantError, err)
			}

			if response != tt.wantResponse {
				t.Errorf("[%s] Resposta esperada: %v, recebida: %v\n", tt.name, tt.wantResponse, response)
			}
		})
	}
}

func TestDeleteContact(t *testing.T) {
    tests := []struct {
        name        string
        userID      int64
        contactID   int64
        mockFunc    func(ctx context.Context, userID, contactID int64) (bool, error)
        wantResp    bool
        wantErr     bool
    }{
        {
            name:       "Contato deleteado",
            userID:     1,
            contactID:  1,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (bool, error) {
                return true, nil
            },
            wantResp: true,
            wantErr: false,
        },
        {
            name:       "Contato não encontrado",
            userID:     1,
            contactID:  2,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (bool, error) {
                return false, nil
            },
            wantResp: false,
            wantErr:  false,
       },
        {
            name:       "Contact ID vazio",
            userID:     1,
            contactID:  internal.ZERO,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (bool, error) {
                return false, ErrContactIDEmpty
            },
            wantResp: false,
            wantErr:  true,
       },
        {
            name:       "UserID vazio",
            userID:     internal.ZERO,
            contactID:  internal.ZERO,
            mockFunc:   func(ctx context.Context, userID, contactID int64) (bool, error) {
                return false, ErrUserIDEmpty
            },
            wantResp: false,
            wantErr:  true,
       },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockContactRepository{
                DeleteContactFunc: tt.mockFunc,
            }
            service := NewContactService(mockRepo)

            response, err := service.DeleteContact(context.Background(), tt.userID, tt.contactID)
            if (err != nil) != tt.wantErr {
                t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v", tt.name, tt.wantErr, err)
            }
            if response != tt.wantResp {
                t.Errorf("[%s] Resposta esperada: %+v, recebida: %+v", tt.name, tt.wantResp, response)
            }
        })
    }
}

var (
	ErrSaveContactFuncNotImplemented   = errors.New("SaveContactFunc not implemented")
	ErrGetContactFuncNotImplemented    = errors.New("GetContactFunc not implemented")
	ErrListContactsFuncNotImplemented  = errors.New("ListContactsFunc not implemented")
	ErrUpdateContactFuncNotImplemented = errors.New("UpdateContactFunc not implemented")
	ErrDeleteContactFuncNotImplemented = errors.New("DeleteContactFunc not implemented")

    ErrUserIDEmpty = errors.New("Err user ID is empty")
    ErrContactIDEmpty = errors.New("Err contact ID is empty")
    ErrInvalidContact = errors.New("Error missing parameters")
    ErrInvalidContactParameters = errors.New("Error invalid  parameters")
    ErrContactNotFound = errors.New("Contact not found")
)
