package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "io"
    "path/filepath"
    "time"

    "github.com/joho/godotenv"
    "github.com/amarantec/move-easy/internal/db"
    "github.com/amarantec/move-easy/internal"
    "github.com/amarantec/move-easy/internal/middleware"
    "github.com/amarantec/move-easy/internal/handlers/routes"
)

func main() {
    loadEnv()
    setupLogger()

    // Wait 10 seconds for podman:postgres start
    time.Sleep(10 * time.Second)

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
    loggedMux := middleware.LoggerMiddleware(mux)

    server := &http.Server {
        Addr:       serverPort,
        Handler:    loggedMux,
    }

    fmt.Printf("Server listen on: localhost%s\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}

func findEnv(path string) (string, error) {
    filePath := filepath.Join(path, ".env")
    if _, err := os.Stat(filePath); err == nil {
        return filePath, nil
    }

    files, err := os.ReadDir(path)
    if err != nil {
        return internal.EMPTY, err
    }
    
    for _, file := range files {
        if file.IsDir() {
            subDirPath := filepath.Join(path, file.Name())
            subDirFilePath, err := findEnv(subDirPath)
            if err == nil {
                return subDirFilePath, nil
            }
        }
    }

    return internal.EMPTY, fmt.Errorf("file .env not found in %s or anywhere\n", path)
}

func loadEnv() {
    path, err := os.Getwd()
    if err != nil {
        log.Fatal("error getting actual dir")
    }

    appPath := filepath.Dir(filepath.Dir(path))
    pathFile, err := findEnv(appPath)
    if err != nil {
        log.Fatal(".env no found")
    }
    godotenv.Load(pathFile)
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

func setupLogger() {
    logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("erro ao abrir o arquivo de log: %v\n", err)
    }

    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
