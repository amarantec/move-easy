package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
    "time"

	"github.com/amarantec/move-easy/internal/utils"
	"github.com/amarantec/move-easy/internal/db"
	"github.com/amarantec/move-easy/internal/handlers/routes"
	"github.com/amarantec/move-easy/internal/middleware"
)

func main() {
	utils.LoadEnv()
	setupLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

	connectionString, err := utils.BuildConnectionString()
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
	}

	fmt.Printf("Server listen on: http://localhost%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func setupLogger() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("erro ao abrir o arquivo de log: \n", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
