package address

import (
    "context"
    "errors"
    "testing"
    "github.com/amarantec/move-easy/internal"
)

type mockAddressRepository struct {
    GetAddressFunc      func(ctx context.Context, userID int64) (internal.Address, error)
    AddOrUpdateFunc     func(ctx context.Context, address internal.Address) (int64, error)
}

func (m *mockAddressRepository) GetAddress(ctx context.Context, userID int64) (internal.Address, error) {
    if m.GetAddressFunc != nil {
        return m.GetAddressFunc(ctx, userID)
    }
    return internal.Address{}, ErrGetAddressNotImplemented
}

func (m *mockAddressRepository) AddOrUpdateAddress(ctx context.Context, address internal.Address) (int64, error) {
    if m.AddOrUpdateFunc != nil {
        return m.AddOrUpdateFunc(ctx, address)
    }
    return internal.ZERO, ErrAddOrUpdateNotImplemented
}

func TestGetAddress(t *testing.T) {
    tests := []struct {
        name      string
        userID    int64
        mockFunc  func(ctx context.Context, userID int64) (internal.Address, error)
        want      internal.Address
        wantErr   bool
    }{
        {
            name: "Endereço encontrado",
            userID: 1,
            mockFunc: func(ctx context.Context, userID int64) (internal.Address, error) {
                return internal.Address{
                    ID:           10,
                    UserID:       1,
                    Street:       "Rua Exemplo",
                    Number:       "123",
                    CEP:          "12345678",
                    Neighborhood: "Centro",
                    City:         "São Paulo",
                    State:        "SP",
                }, nil
            },
            want: internal.Address{
                ID:           10,
                UserID:       1,
                Street:       "Rua Exemplo",
                Number:       "123",
                CEP:          "12345678",
                Neighborhood: "Centro",
                City:         "São Paulo",
                State:        "SP",
            },
            wantErr: false,
        },
        {
            name: "Usuário sem endereço cadastrado",
            userID: 2,
            mockFunc: func(ctx context.Context, userID int64) (internal.Address, error) {
                return internal.Address{}, ErrAddressNotFound
            },
            want:    internal.Address{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockAddressRepository{
                GetAddressFunc: tt.mockFunc,
            }
            service := NewAddressService(mockRepo)

            got, err := service.GetAddress(context.Background(), tt.userID)
            if (err != nil) != tt.wantErr {
                t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v", tt.name, tt.wantErr, err)
            }
            if got != tt.want {
                t.Errorf("[%s] Endereço esperado: %+v, recebido: %+v", tt.name, tt.want, got)
            }
        })
    }
}

func TestAddOrUpdateAddress(t *testing.T) {
    tests := []struct {
        name     string
        input    internal.Address
        mockFunc func(ctx context.Context, address internal.Address) (int64, error)
        wantID   int64
        wantErr  bool
    }{
        {
            name: "Endereço adicionado com sucesso",
            input: internal.Address{
                UserID:       1,
                Street:       "Rua Nova",
                Number:       "456",
                CEP:          "98765432",
                Neighborhood: "Bairro Novo",
                City:         "Rio de Janeiro",
                State:        "RJ",
            },
            mockFunc: func(ctx context.Context, address internal.Address) (int64, error) {
                return 20, nil
            },
            wantID:  20,
            wantErr: false,
        },
        {
            name: "Erro ao adicionar endereço",
            input: internal.Address{
                UserID:       2,
                Street:       "",
                Number:       "",
                CEP:          "",
                Neighborhood: "",
                City:         "",
                State:        "",
            },
            mockFunc: func(ctx context.Context, address internal.Address) (int64, error) {
                return internal.ZERO, ErrInvalidAddress
            },
            wantID:  internal.ZERO,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockAddressRepository{
                AddOrUpdateFunc: tt.mockFunc,
            }
            service := NewAddressService(mockRepo)

            id, err := service.AddOrUpdateAddress(context.Background(), tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v", tt.name, tt.wantErr, err)
            }
            if id != tt.wantID {
                t.Errorf("[%s] ID esperado: %d, recebido: %d", tt.name, tt.wantID, id)
            }
        })
    }
}

var (
    ErrGetAddressNotImplemented = errors.New("GetAddress function not implemented")
    ErrAddOrUpdateNotImplemented = errors.New("AddOrUpdate function not implemented")
    ErrAddressNotFound           = errors.New("Address not found")
    ErrInvalidAddress            = errors.New("Invalid address details")
)

