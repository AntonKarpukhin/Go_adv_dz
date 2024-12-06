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

	verify.NewVerifierHandler(router, file.NewJsFile("data.json"))

	server := http.Server{
		Addr:    ":8083",
		Handler: router,
	}

	fmt.Println("Server is listening 8083")
	server.ListenAndServe()
}
