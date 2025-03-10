package user

import (
    "context"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/utils"
)

type IUserService interface {
    Register (ctx context.Context, user internal.UserRegister) (int64, error)
    ValidateCredentials(ctx context.Context, user internal.UserLogin) (string, error)
}

type userService struct {
    userRepository IUserRepository
}

func NewUserService(repository IUserRepository) IUserService {
    return &userService{userRepository: repository}
}

func (s *userService) Register (ctx context.Context, user internal.UserRegister) (int64, error) {
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        return internal.ZERO, err
    }

    user.Password = hashedPassword

    response, err := s.userRepository.Register(ctx, user)
    if err != nil {
        return internal.ZERO, err
    }

    return response, nil
}

func (s *userService) ValidateCredentials(ctx context.Context, user internal.UserLogin) (string, error) {
    userDb, err := s.userRepository.ValidateCredentials(ctx, user)
    if err != nil {
        return internal.EMPTY, err
    }

    passwordIsValid :=
        utils.CheckPasswordHash(user.Password, userDb.Password)
    if !passwordIsValid {
        return internal.EMPTY, nil
    }

    token, err := utils.GenerateToken(userDb.Email, userDb.ID)
    if err != nil {
        return internal.EMPTY, err
     }

     return token, nil
}
