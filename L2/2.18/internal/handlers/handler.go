package handler

import (
	"net/http"

	"calendar/internal/usecases"
)

type Handler struct {
	calendar *usecases.Calendar
}

func NewHandler(c *usecases.Calendar) *Handler {
	return &Handler{calendar: c}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request)
