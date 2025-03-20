package handlers

import (
    "net/http"
    "encoding/json"
    "time"
    "context"
    "strconv"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/contact"
    "github.com/amarantec/move-easy/internal/middleware"
)

type ContactHandler struct {
    service     contact.IContactService
}

func NewContactHandler(service contact.IContactService) *ContactHandler {
    return &ContactHandler{service: service}
}

func (h *ContactHandler) SaveContact(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    userID := r.Context().Value(middleware.UserIDKey).(int64)

    var contact internal.Contact
    contact.UserID = userID

    if err :=
        json.NewDecoder(r.Body).Decode(&contact); err != nil {
            http.Error(w,
                "could not decode this request, error: " + err.Error(),
                http.StatusBadRequest)
        return
    }

    response, err := h.service.SaveContact(ctxTimeout, contact)
    if err != nil {
        http.Error(w,
            "could not save this contact, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response,
    })
}

func (h *ContactHandler) GetContact(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    userID := r.Context().Value(middleware.UserIDKey).(int64)

    contactID, err := strconv.ParseInt(r.PathValue("contactID"), 10, 64)
    if err != nil {
        http.Error(w,
            "invalid parameter, error: " + err.Error(),
            http.StatusBadRequest)
        return
    }

    response, err := h.service.GetContact(ctxTimeout, userID, contactID)
    if err != nil {
        http.Error(w,
           "could not get this contact, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response,
    })
}

func (h *ContactHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    userID := r.Context().Value(middleware.UserIDKey).(int64)

    response, err := h.service.ListContacts(ctxTimeout, userID)
    if err != nil {
        http.Error(w,
           "could not get the contact list, error: " + err.Error(),
            http.StatusInternalServerError)
        return

    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response,
    })
}

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    userID := r.Context().Value(middleware.UserIDKey).(int64)

    var contact internal.Contact
    contact.UserID = userID

    if err :=
        json.NewDecoder(r.Body).Decode(&contact); err != nil {
            http.Error(w,
                "could not decode this request, error: " + err.Error(),
                http.StatusBadRequest)
        return
    }

    response, err := h.service.UpdateContact(ctxTimeout, contact)
    if err != nil {
        http.Error(w,
            "could not update this contact, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNoContent)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response,
    })
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    userID := r.Context().Value(middleware.UserIDKey).(int64)

    contactID, err := strconv.ParseInt(r.PathValue("contactID"), 10, 64)
    if err != nil {
        http.Error(w,
            "invalid parameter, error: " + err.Error(),
            http.StatusBadRequest)
        return
    }

    response, err := h.service.DeleteContact(ctxTimeout, userID, contactID)
    if err != nil {
        http.Error(w,
           "could not delete this contact, error: " + err.Error(),
            http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNoContent)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response,
    })
}
