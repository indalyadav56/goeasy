package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/test/authtest-fixed/internal/product/interface/http/v1/handlers"
)

func SetupProductRoutes(r chi.Router, h handlers.ProductHandler) {
	r.Route("/v1/product", func(r chi.Router) {
		r.Get("/", h.GetProduct)
		r.Post("/", h.CreateProduct)
		r.Put("/{id}", h.UpdateProduct)
		r.Delete("/{id}", h.DeleteProduct)
	})
}
