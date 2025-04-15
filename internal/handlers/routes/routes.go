package routes

import (
	"net/http"

	"github.com/amarantec/move-easy/internal/address"
	"github.com/amarantec/move-easy/internal/bus"
	"github.com/amarantec/move-easy/internal/contact"
	"github.com/amarantec/move-easy/internal/handlers"
	"github.com/amarantec/move-easy/internal/sharedVehicle"
	"github.com/amarantec/move-easy/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetRoutes(conn *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	/*
	   Address Dependency Injection
	*/
	addrRepository := address.NewAddressRepository(conn)
	addrService := address.NewAddressService(addrRepository)
	addrHandler := handlers.NewAddressHandler(addrService)

	/*
	   User Dependency Injection
	*/

	userRepository := user.NewUserRepository(conn)
	userService := user.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	/*
	   Contact Dependency Injection
	*/

	contactRepository := contact.NewContactRepository(conn)
	contactService := contact.NewContactService(contactRepository)
	contactHandler := handlers.NewContactHandler(contactService)

	/*
		Shared Vehicle Dependency Injection
	*/

	sharedVehicleRepository := sharedVehicle.NewSharedVehicleRepository(conn)
	sharedVehicleService := sharedVehicle.NewSharedVehicleService(sharedVehicleRepository)
	sharedVehicleHandler := handlers.NewSharedVehicleHandler(sharedVehicleService)

	/*
		Bus Dependency Injection
	*/
	busRepository := bus.NewBusRepository(conn)
	busService := bus.NewBusService(busRepository)
	busHandler := handlers.NewBusHandler(busService)

	/*
	   Routes
	*/

	mux.Handle("/user/", http.StripPrefix("/user", userRoutes(userHandler)))
	mux.Handle("/address/", http.StripPrefix("/address", addressRoutes(addrHandler)))
	mux.Handle("/contact/", http.StripPrefix("/contact", contactRoutes(contactHandler)))
	mux.Handle("/shared-vehicle/", http.StripPrefix("/shared-vehicle", sharedVehicleRoutes(sharedVehicleHandler)))
	mux.Handle("/bus/", http.StripPrefix("/bus", busRoutes(busHandler)))
	return mux
}
