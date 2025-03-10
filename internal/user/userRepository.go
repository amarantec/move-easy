package user

import (
    "context"
    "github.com/amarantec/move-easy/internal"
    "github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
    Register(ctx context.Context, user internal.UserRegister) (int64, error)
    ValidateCredentials(ctx context.Context, user internal.UserLogin) (internal.UserLogin, error)
}

type userRepository struct {
    Conn    *pgxpool.Pool
}

func NewUserRepository(connection *pgxpool.Pool) IUserRepository {
    return &userRepository{Conn: connection}
}

func (r *userRepository) Register(ctx context.Context, user internal.UserRegister) (int64, error) {
    var userID int64
    err :=
        r.Conn.QueryRow(
            ctx,
            `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id;`,
            user.Email, user.Password).Scan(&userID)
    if err != nil {
        return internal.ZERO, err
    }

    return userID, nil
}

func (r *userRepository) ValidateCredentials(ctx context.Context, user internal.UserLogin) (internal.UserLogin, error) {
    err :=
        r.Conn.QueryRow(
            ctx,
            `SELECT id, password FROM users WHERE email = $1;`, user.Email).Scan(&user.ID, &user.Password)

    if err != nil {
        return internal.UserLogin{}, err
    }

    return user, nil
}
