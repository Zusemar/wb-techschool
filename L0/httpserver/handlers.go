package httpserver

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"wb-techschool/L0/internal/domain"
	"wb-techschool/L0/internal/usecases"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Repo  usecases.OrderRepository
	Cache *usecases.OrderCache
}

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "order_uid")
	if id == "" {
		http.Error(w, "order_uid is required", http.StatusBadRequest)
		return
	}

	if s.Cache != nil {
		if o, ok := s.Cache.Get(id); ok {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(o)
			return
		}
	}

	order, err := s.Repo.GetOrderById(r.Context(), id)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}

func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o domain.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if o.Order_uid == "" {
		http.Error(w, "order_uid is required", http.StatusBadRequest)
		return
	}
	if err := s.Repo.CreateOrder(r.Context(), &o); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if s.Cache != nil {
		s.Cache.Set(o)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(o)
}
