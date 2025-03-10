package handlers

import (
    "net/http"
    "encoding/json"
    "context"
    "time"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/user"
)

type UserHandler struct {
    service user.IUserService
}

func NewUserHandler(service user.IUserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    var user internal.UserRegister

    if err :=
        json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w,
                "could not decode this request, error: " + err.Error(),
                http.StatusBadRequest)
        return
    }

    response, err := h.service.Register(ctxTimeout, user)
    if err != nil {
        http.Error(w,
            "could not register this user, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response,
    })
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    var user internal.UserLogin

    if err :=
        json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w,
                "could not decode this request, error: " + err.Error(),
                http.StatusBadRequest)
        return
    }

    response, err := h.service.ValidateCredentials(ctxTimeout, user)
    if err != nil {
        http.Error(w,
            "could not validate this credentials, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "token": response,
    })
}
