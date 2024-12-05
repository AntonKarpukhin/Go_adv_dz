package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	NewRandomAPI(router)

	server := http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	fmt.Println("Start server")
	server.ListenAndServe()
}
