package product

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"orderApi/pkg/request"
	"orderApi/pkg/response"
	"strconv"
)

type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, productDeps ProductHandler) {
	handler := &ProductHandler{
		ProductRepository: productDeps.ProductRepository,
	}

	router.HandleFunc("POST /product", handler.Create())
	router.HandleFunc("PATCH /product/{id}", handler.Update())
	router.HandleFunc("DELETE /product/{id}", handler.Delete())
	router.HandleFunc("GET /product/{id}", handler.FindId())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.BodyDecode[ProductCreateRequest](&w, r)
		if err != nil {
			return
		}
		product := NewProduct(body)

		existedProduct, _ := handler.ProductRepository.GetByName(product.Name)
		if existedProduct != nil {
			fmt.Println("product existed:", product.Name)
			return
		}

		createdProduct, errCreate := handler.ProductRepository.Create(product)
		if errCreate != nil {
			http.Error(w, errCreate.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, createdProduct, http.StatusOK)
	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.BodyDecode[ProductUpdateRequest](&w, r)
		if err != nil {
			return
		}

		isString := r.PathValue("id")
		id, errParseUint := strconv.ParseUint(isString, 10, 32)
		if errParseUint != nil {
			http.Error(w, errParseUint.Error(), http.StatusBadRequest)
			return
		}

		product, errUpdate := handler.ProductRepository.Update(&Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
		})

		if errUpdate != nil {
			http.Error(w, errUpdate.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.ProductRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.ProductRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Json(w, "Запись удалена", http.StatusOK)
	}
}

func (handler *ProductHandler) FindId() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		product, errFind := handler.ProductRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, errFind.Error(), http.StatusNotFound)
			return
		}
		response.Json(w, product, http.StatusOK)
	}
}
