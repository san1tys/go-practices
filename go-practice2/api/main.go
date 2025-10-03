package main

import (
	"fmt"
	"log"
	"net/http"

	"go-practice2/internal/handlers"
	"go-practice2/internal/middleware"
)

func main() {
	userHandler := handlers.NewUserHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /user", userHandler.GetUser)
	mux.HandleFunc("POST /user", userHandler.CreateUser)

	handler := middleware.AuthMiddleware(mux)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Test with: curl -i -H \"X-API-Key: secret123\" \"http://localhost%s/user?id=42\"\n", port)

	log.Fatal(http.ListenAndServe(port, handler))
}
