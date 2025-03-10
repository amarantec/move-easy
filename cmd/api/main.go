package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/amarantec/move-easy/internal/db"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/handlers/routes"
)

func main() {
    loadEnv()
    ctx := context.Background()
    serverPort := ":8080"

    connectionString, err := buildConnectionString()
    if err != nil {
        log.Fatal(err)
    }

    Conn, err := db.OpenConnection(ctx, connectionString)
    if err != nil {
        panic(err)
    }

    defer Conn.Close()

    mux := routes.SetRoutes(Conn)

    server := &http.Server {
        Addr:       serverPort,
        Handler:    mux,
    }

    fmt.Printf("Server listen on: localhost%s\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}

func loadEnv() {
    err := godotenv.Load("../../config/.env")
    if err != nil {
        log.Fatal("error loading .env file")
    }
}

func buildConnectionString() (string, error) {
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")
    dbPort := os.Getenv("DB_PORT")

    if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
        return internal.EMPTY, fmt.Errorf("one or more environment variables are not set")
    }

    connectionString := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
                        dbHost, dbPort, dbUser, dbPassword, dbName)

    return connectionString, nil
}
