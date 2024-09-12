package v1

import (
	"avitoTech/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(services *service.Services) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Pong"))
		})

		r.Route("/tenders", func(r chi.Router) {
			newTenderRoutes(r, services.Tender)
		})
	})

	//r.Route("/test", func(r chi.Router) {
	//	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("Pong"))
	//	})
	//})

	return r
}
