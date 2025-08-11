package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(s *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", s.Healthz)
	r.Get("/order/{order_uid}", s.GetOrder)
	r.Post("/order", s.CreateOrder)

	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	return r
}
