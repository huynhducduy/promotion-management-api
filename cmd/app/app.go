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
	"promotion-management-api/internal/member"
	"promotion-management-api/internal/order"
	"promotion-management-api/internal/product"
	"promotion-management-api/internal/promotion"
	"promotion-management-api/internal/store"
)

func Run() error {

	config.ReadConfig()

	db.OpenConnection()
	auth.InitFirebase()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(c.Handler)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", auth.GetPwd)
		r.Post("/login", auth.Login)

		r.Route("/", func(r chi.Router) {

			r.Use(auth.AuthenticationMiddleware)

			r.Get("/me", employee.RouterMe)

			r.Route("/promotion", func(r chi.Router) {
				r.Get("/", promotion.RouterList)
				r.Post("/", promotion.RouterCreate)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", promotion.RouterRead)
					r.Put("/", promotion.RouterUpdate)
					r.Delete("/", promotion.RouterDelete)
				})

				r.Get("/applicable", promotion.RouterApplicable)
			})

			r.Route("/order", func(r chi.Router) {
				r.Get("/", order.RouterList)
				r.Post("/", order.RouterCreate)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", order.RouterRead)
					r.Put("/", order.RouterUpdate)
					r.Delete("/", order.RouterDelete)
				})
			})

			r.Route("/product", func(r chi.Router) {
				r.Get("/", product.RouterList)
				//r.Post("/", product.RouterCreate)

				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", product.RouterRead)
				//	r.Put("/", product.RouterUpdate)
				//	r.Delete("/", product.RouterDelete)
				//})
			})

			r.Route("/store", func(r chi.Router) {
				r.Get("/", store.RouterList)
				//r.Post("/", store.RouterCreate)

				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", store.RouterRead)
				//	r.Put("/", store.RouterUpdate)
				//	r.Delete("/", store.RouterDelete)
				//})
			})

			r.Route("/employee", func(r chi.Router) {
				r.Get("/", employee.RouterList)
				//r.Post("/", employee.RouterCreate)

				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", employee.RouterRead)
				//	r.Put("/", employee.RouterUpdate)
				//	r.Delete("/", employee.RouterDelete)
				//})
			})

			r.Route("/member", func(r chi.Router) {
				r.Get("/", member.RouterList)
				//r.Post("/", member.RouterCreate)

				//r.Route("/{id}", func(r chi.Router) {
				//	r.Get("/", member.RouterRead)
				//	r.Put("/", member.RouterUpdate)
				//	r.Delete("/", member.RouterDelete)
				//})
			})
		})
	})

	log.Printf("Running at port 80")
	return http.ListenAndServe(":80", r)
}
