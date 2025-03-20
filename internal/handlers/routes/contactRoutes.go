package routes

import (
    "github.com/amarantec/move-easy/internal/handlers"
    "github.com/amarantec/move-easy/internal/middleware"
    "net/http"
)

func contactRoutes(handler *handlers.ContactHandler) *http.ServeMux {
    contactMux := http.NewServeMux()

    contactMux.HandleFunc("/save-contact", middleware.Authenticate(handler.SaveContact))
    contactMux.HandleFunc("/get-contact/{contactID}", middleware.Authenticate(handler.GetContact))
    contactMux.HandleFunc("/list-contacts", middleware.Authenticate(handler.ListContacts))
    contactMux.HandleFunc("/update-contact", middleware.Authenticate(handler.UpdateContact))
    contactMux.HandleFunc("/delete-contact/{contactID}", middleware.Authenticate(handler.DeleteContact))

    return contactMux
}
