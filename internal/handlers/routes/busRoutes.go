package routes

import (
	"net/http"

	"github.com/amarantec/move-easy/internal/handlers"
)

func busRoutes(handler *handlers.BusHandler) *http.ServeMux {
	busMux := http.NewServeMux()

	busMux.HandleFunc("/insert-new-bus-line", handler.InsertNewBusLine)
	busMux.HandleFunc("/insert-bus-stop", handler.InsertBusStop)
	busMux.HandleFunc("/get-bus-line/{busLineID}", handler.GetBusLine)
	busMux.HandleFunc("/get-bus-stop/{busStopID}", handler.GetBusStop)

	return busMux
}
