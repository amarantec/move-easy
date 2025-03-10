package utils

import (
    "log"
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/amarantec/move-easy/internal"
)

const secret = "secret"

func GenerateToken(email string, userID int64) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
        "userID": userID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })

    return token.SignedString([]byte(secret))
}

func VerifyToken(token string) (int64, error) {
    parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        _, ok := t.Method.(*jwt.SigningMethodHMAC)
        if !ok {
            return internal.ZERO, ErrUnexpectedSigningMethod
        }

        return []byte(secret), nil
    })

    if err != nil {
        log.Print("Could not parse this token.\n")
        return internal.ZERO, ErrCouldNotParseToken
    }

    tokenIsValid := parsedToken.Valid
    if !tokenIsValid {
        return internal.ZERO, ErrInvalidToken
    }

    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok {
        return internal.ZERO, ErrInvalidTokenClaims
    }

    userID := int64(claims["userID"].(float64))
    return userID, nil
}

var ErrUnexpectedSigningMethod = errors.New("Unexpected signing method")
var ErrCouldNotParseToken = errors.New("Could not parse token")
var ErrInvalidToken = errors.New("Invalid Token")
var ErrInvalidTokenClaims = errors.New("Invalid token claims")
