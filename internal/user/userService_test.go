package user

import (
    "context"
    "errors"
    "testing"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/utils"
)

type mockUserRepository struct {
    RegisterFunc func (ctx context.Context, user internal.UserRegister) (int64, error)
    ValidateCredentialsFunc func (ctx context.Context, user internal.UserLogin) (internal.UserLogin, error) 
}

func (m *mockUserRepository) Register(ctx context.Context, user internal.UserRegister) (int64, error) {
    if m.RegisterFunc != nil {
        return m.RegisterFunc(ctx, user)
    }

    return internal.ZERO, ErrRegisterFuncNotImplemented
}


func (m *mockUserRepository) ValidateCredentials(ctx context.Context, user internal.UserLogin) (internal.UserLogin, error) {
    if m.ValidateCredentialsFunc != nil {
        return m.ValidateCredentialsFunc(ctx, user)
    }
    return internal.UserLogin{}, ErrValidateCredentialsFuncNotImplemented
}

func TestRegister(t *testing.T) {
    tests := []struct {
        name        string
        input       internal.UserRegister
        mockFunc    func(ctx context.Context, user internal.UserRegister) (int64, error)
        wantID      int64
        wantError   bool
    }{
        {
            name: "Cadastro bem-sucedido",
            input: internal.UserRegister{
                Email: "valid@example.com",
                Password: "StrongPass123",
            },
            mockFunc: func(ctx context.Context, user internal.UserRegister) (int64, error) {
                return 1, nil
            },
            wantID: 1,
            wantError: false,
        },
        {
            name: "Erro ao tentar cadastrar usuário com e-mail vazio",
            input: internal.UserRegister{
                Email: "",
                Password: "StrongPass123",
            },
            mockFunc: func(ctx context.Context, user internal.UserRegister) (int64, error) {
                return internal.ZERO, ErrEmailEmpty
            },
            wantID: internal.ZERO,
            wantError: true,
        },
        {
            name: "Erro ao tentar cadastrar usuário com senha vazia",
            input: internal.UserRegister{
                Email: "valid@example.com",
                Password: "",
            },
            mockFunc: func(ctx context.Context, user internal.UserRegister) (int64, error) {
                return internal.ZERO, ErrPasswordEmpty
            },
            wantID: internal.ZERO,
            wantError: true,
        },
        {
            name: "Erro simulado do banco de dados",
            input: internal.UserRegister{
                Email: "db@error.com",
                Password: "SomePassword",
            },
            mockFunc: func(ctx context.Context, user internal.UserRegister) (int64, error) {
                return internal.ZERO, ErrDatabaseError
            },
            wantID: internal.ZERO,
            wantError: true,
        },
    }

    // Executando cada teste
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockUserRepository {
                RegisterFunc: func(ctx context.Context, user internal.UserRegister) (int64, error) {
                    if user.Password != "" {
                        hashedPassword, err := utils.HashPassword(user.Password)
                        if err != nil {
                            return internal.ZERO, err
                        }

                        user.Password = hashedPassword
                    }

                    return tt.mockFunc(ctx, user)
                },
            }

            service := NewUserService(mockRepo)

            id, err := service.Register(context.Background(), tt.input)
            if (err != nil) != tt.wantError {
                t.Errorf("[%s] Esperava erro: %v, recebeu erro: %v", tt.name, tt.wantError, err)
            }

            if id != tt.wantID {
                t.Errorf("[%s] ID esperado: %d, recebido: %d", tt.name, tt.wantID, id)
            }
        })
    }
}

var (
    ErrRegisterFuncNotImplemented = errors.New("RegisterFunc not implemented")
    ErrValidateCredentialsFuncNotImplemented = errors.New("ValidateCredentialsFunc not implemented")
    ErrEmailEmpty = errors.New("Email cannot be empty")
    ErrPasswordEmpty = errors.New("Password cannot be empty")
    ErrDatabaseError = errors.New("Database error")
)
