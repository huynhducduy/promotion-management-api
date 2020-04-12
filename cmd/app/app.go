package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"promotion-management-api/internal/auth"
	"promotion-management-api/internal/config"
	"promotion-management-api/internal/db"
	"promotion-management-api/internal/employee"
	"promotion-management-api/internal/order"
	"promotion-management-api/internal/product"
	"promotion-management-api/internal/promotion"
	"promotion-management-api/internal/store"
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
		r.Get("/", auth.GetPwd)

		r.Route("/", func(r chi.Router) {

			r.Use(auth.AuthenticationgMiddleware)

			r.Route("/promotion", func(r chi.Router) {
				r.Get("/", promotion.RouterList)
				r.Post("/", promotion.RouterCreate)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", promotion.RouterRead)
					r.Post("/", promotion.RouterUpdate)
					r.Delete("/", promotion.RouterDelete)
				})
			})

			r.Route("/order", func(r chi.Router) {
				r.Get("/", order.RouterList)
				//r.Post("/", order.RouterCreate)
				//
				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", order.RouterRead)
				//	r.Post("/", order.RouterUpdate)
				//	r.RouterDelete("/", order.RouterDelete)
				//})
			})

			r.Route("/product", func(r chi.Router) {
				r.Get("/", product.RouterList)
				//r.Post("/", product.RouterCreate)
				//
				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", product.RouterRead)
				//	r.Post("/", product.RouterUpdate)
				//	r.RouterDelete("/", product.RouterDelete)
				//})
			})

			r.Route("/store", func(r chi.Router) {
				r.Get("/", store.RouterList)
				//r.Post("/", store.RouterCreate)
				//
				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", store.RouterRead)
				//	r.Post("/", store.RouterUpdate)
				//	r.RouterDelete("/", store.RouterDelete)
				//})
			})

			r.Route("/employee", func(r chi.Router) {
				r.Post("/login", auth.Login)

				r.Get("/", employee.RouterList)
				//r.Post("/", employee.RouterCreate)
				//
				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", employee.RouterRead)
				//	r.Post("/", employee.RouterUpdate)
				//	r.RouterDelete("/", employee.RouterDelete)
				//})
			})

			r.Route("/member", func(r chi.Router) {
				//r.Get("/", member.RouterList)
				//r.Post("/", member.RouterCreate)
				//
				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", member.RouterRead)
				//	r.Post("/", member.RouterUpdate)
				//	r.RouterDelete("/", member.RouterDelete)
				//})
			})
		})
	})

	log.Printf("Running at port 80")
	return http.ListenAndServe(":80", r)
}
