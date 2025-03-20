package routes

import (
	"github.com/amarantec/move-easy/internal/handlers"
	"net/http"
)

func userRoutes(handler *handlers.UserHandler) *http.ServeMux {
	userMux := http.NewServeMux()

	userMux.HandleFunc("/register", handler.Register)
	userMux.HandleFunc("/login", handler.Login)

	return userMux
}
