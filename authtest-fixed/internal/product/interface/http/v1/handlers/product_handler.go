package handlers

import (
	"net/http"
	"github.com/test/authtest-fixed/internal/product/application"
)

type ProductHandler interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	ListProducts(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	service application.ProductService
}

func NewProductHandler(service application.ProductService) ProductHandler {
	return &productHandler{
		service: service,
	}
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	
}


func (h *productHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	
}


func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	
}


func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	
}

func (h *productHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	
}