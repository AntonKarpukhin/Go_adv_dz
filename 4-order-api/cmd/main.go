package main

import (
	"fmt"
	"net/http"
	"orderApi/configs"
	"orderApi/iternal/auth"
	"orderApi/iternal/product"
	"orderApi/iternal/user"
	"orderApi/pkg/db"
	"orderApi/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	//Repositories
	productRepository := product.NewProductRepository(database)
	userRepository := user.NewUserRepository(database)

	//Services
	authServices := auth.NewAuthService(userRepository)

	//Handler
	product.NewProductHandler(router, product.ProductHandler{
		ProductRepository: productRepository,
	})

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authServices,
	})

	//Middlewares
	stack := middleware.Chain(middleware.Logrus)

	server := http.Server{
		Addr:    ":8086",
		Handler: stack(router),
	}
	fmt.Println("Server is listening 8086")
	server.ListenAndServe()
}
