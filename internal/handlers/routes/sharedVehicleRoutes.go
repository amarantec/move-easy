package routes

import (
	"net/http"

	"github.com/amarantec/move-easy/internal/handlers"
	"github.com/amarantec/move-easy/internal/middleware"
)

func sharedVehicleRoutes(handler *handlers.SharedVehicleHandler) *http.ServeMux {
	sharedVehicleMux := http.NewServeMux()

	sharedVehicleMux.HandleFunc("/insert-shared-vehicle", middleware.Authenticate(handler.InsertSharedVehicle))
	sharedVehicleMux.HandleFunc("/get-shared-vehicle/{vehicleID}", handler.GetSharedVehicle)
	sharedVehicleMux.HandleFunc("/list-shared-vehicles", handler.ListAllSharedVehicles)
	sharedVehicleMux.HandleFunc("/update-shared-vehicle-location", middleware.Authenticate(handler.UpdateSharedVehicleLocation))

	return sharedVehicleMux
}
