package routes

import (
	"github.com/amarantec/move-easy/internal/handlers"
	"github.com/amarantec/move-easy/internal/middleware"
	"net/http"
)

func addressRoutes(handler *handlers.AddressHandler) *http.ServeMux {
	addrMux := http.NewServeMux()

	addrMux.HandleFunc("/get-address", middleware.Authenticate(handler.GetAddress))
	addrMux.HandleFunc("/save-address", middleware.Authenticate(handler.AddOrUpdateAddress))

	return addrMux
}
