package routes

import(
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/amarantec/move-easy/internal/address"
    "github.com/amarantec/move-easy/internal/user"
    "github.com/amarantec/move-easy/internal/handlers"
    "github.com/amarantec/move-easy/internal/contact"
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
        Contact Dependency Injection
    */

    contactRepository := contact.NewContactRepository(conn)
    contactService := contact.NewContactService(contactRepository)
    contactHandler := handlers.NewContactHandler(contactService)

    /*
        Routes
    */

    mux.Handle("/user/", http.StripPrefix("/user", userRoutes(userHandler)))
    mux.Handle("/address/", http.StripPrefix("/address", addressRoutes(addrHandler)))
    mux.Handle("/contact/", http.StripPrefix("/contact", contactRoutes(contactHandler)))
    return mux
}
