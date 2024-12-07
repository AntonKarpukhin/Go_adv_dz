package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"validation/internal/verify"
	"validation/pkg/file"
)

func main() {
	godotenv.Load()
	router := http.NewServeMux()

	// Создаем обработчик
	verifierHandler := NewVerifierHandler(router, db)

	// Останавливаем горутину перед завершением приложения
	defer verifierHandler.Stop()

	server := http.Server{
		Addr:    ":8083",
		Handler: router,
	}

	fmt.Println("Server is listening 8083")
	server.ListenAndServe()
}
