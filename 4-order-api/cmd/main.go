package main

import (
	"fmt"
	"net/http"
	"orderApi/configs"
	"orderApi/iternal/product"
	"orderApi/pkg/db"
	"orderApi/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	//Repositories
	productRepository := product.NewProductRepository(database)

	//Handler
	product.NewProductHandler(router, product.ProductHandler{
		ProductRepository: productRepository,
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
