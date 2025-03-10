package handlers

import (
    "net/http"
    "encoding/json"
    "context"
    "time"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/address"
    "github.com/amarantec/move-easy/internal/middleware"
)

type AddressHandler struct {
    service address.IAddressService
}

func NewAddressHandler(service address.IAddressService) *AddressHandler {
    return &AddressHandler{service: service}
}

func (h *AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    userID := r.Context().Value(middleware.UserIDKey).(int64)

    addr, err := h.service.GetAddress(ctxTimeout, userID)
    if err != nil {
        http.Error(w,
            "could not get this address, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": addr,
    })
}

func (h *AddressHandler) AddOrUpdateAddress(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    userID := r.Context().Value(middleware.UserIDKey).(int64)

    var addr internal.Address

    if err :=
        json.NewDecoder(r.Body).Decode(&addr); err != nil {
            http.Error(w,
                "could not decode this request, error: " + err.Error(),
                http.StatusBadRequest)
        return
    }

    addr.UserID = userID
    response, err := h.service.AddOrUpdateAddress(ctxTimeout, addr)
    if err != nil {
        http.Error(w,
            "could not save this address, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response,
    })
}
