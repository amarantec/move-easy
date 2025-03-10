package middleware

import (
    "context"
    "net/http"
    "github.com/amarantec/move-easy/internal/utils"
    "github.com/amarantec/move-easy/internal"
)

type contextKey string
const UserIDKey contextKey = "userID"

func Authenticate (next http.HandlerFunc) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == internal.EMPTY {
            http.Error(w,
                "Token is empty",
                http.StatusUnauthorized)
            return
        }

        userID, err := utils.VerifyToken(token)
        if err != nil {
            http.Error(w,
                err.Error(),
                http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next(w, r.WithContext(ctx))
    }
}
