package middleware

import (
    "context"
    "net/http"
    "github.com/amarantec/move-easy/internal/utils"
)

type contextKey string
const UserIDKey contextKey = "userID"

func Authenticate (next http.HandlerFunc) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            http.Error(w,
                "Token is empty, error: " + err.Error(),
                http.StatusUnauthorized)
            return
        }
    
        token := cookie.Value

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
