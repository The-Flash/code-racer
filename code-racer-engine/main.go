package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/The-Flash/code-racer/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	router := mux.NewRouter()
	var port string
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8000" // default port
	}

	router.
		HandleFunc("/api/v1/exec", handlers.CodeExecHandler).
		Methods("POST").
		Headers("Content-Type", "application/json")

	fmt.Println("Server is running on port", port)
	http.ListenAndServe(":"+port, router)
}
