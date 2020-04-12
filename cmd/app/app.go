package app

import (
	"net/http"
	"promotion-management-api/internal/promotion"

	//"promotion-management-api/pkg/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"promotion-management-api/internal/config"
	"promotion-management-api/internal/db"
)

func Run() error {

	config.ReadConfig()

	db.OpenConnection()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/promotion", func(r chi.Router) {
			r.Get("/", promotion.List)
			r.Post("/", promotion.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", promotion.Read)
				r.Post("/", promotion.Update)
				r.Delete("/", promotion.Delete)
			})
		})
	})

	log.Printf("Running at port 80")
	return http.ListenAndServe(":80", r)
}
