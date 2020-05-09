package router

import (
	"apiservice/handler"
	"apiservice/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func NewRouter(p *handler.Provider) chi.Router {

	r := chi.NewRouter()

	//Cors to handle cross origin requests
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "token"},
		AllowCredentials: true,
		MaxAge:           300,
	})


	r.Use(c.Handler)
	r.Use(middleware.RecoverHandler)
	r.Use(middleware.LoggingHandler)

	r.Get("/ping", p.Ping)

	r.Route("/api/v1", func(sr chi.Router) {
		sr.Post("/write/info", p.InsertNewInfo)
		sr.Get("/read/info", p.ReadInfo)
	})

	return r
}
