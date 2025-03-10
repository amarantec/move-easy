package routes

import(
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/amarantec/move-easy/internal/address"
    "github.com/amarantec/move-easy/internal/user"
    "github.com/amarantec/move-easy/internal/handlers"
    "github.com/amarantec/move-easy/internal/middleware"
    "net/http"
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
        Routes
    */

    mux.HandleFunc("/user/register", userHandler.Register)
    mux.HandleFunc("/user/login", userHandler.Login)


    mux.HandleFunc("/address/get-address", middleware.Authenticate(addrHandler.GetAddress))
    mux.HandleFunc("/address/save-address", middleware.Authenticate(addrHandler.AddOrUpdateAddress))

    return mux
}
