package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"notifier/internal/domain"
	"notifier/internal/usecases"
)

type Handler struct {
	notifier *usecases.Notifier
}

func NewHandler(n *usecases.Notifier) *Handler {
	return &Handler{notifier: n}
}

func (h *Handler) RegisterRoutes(engine *ginext.Engine) {
	engine.Use(ginext.Logger(), ginext.Recovery())

	engine.GET("/", func(c *ginext.Context) {
		c.File("static/index.html")
	})

	api := engine.Group("/notify")
	api.POST("", h.create)
	api.GET("", h.listAll)
	api.GET("/:id", h.getByID)
	api.DELETE("/:id", h.cancel)
}

type createRequest struct {
	Title       string `json:"title" binding:"required"`
	Message     string `json:"message" binding:"required"`
	Channel     string `json:"channel"`
	Recipient   string `json:"recipient"`
	ScheduledAt string `json:"scheduled_at" binding:"required"`
}

func (h *Handler) create(c *ginext.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid scheduled_at, use RFC3339 (e.g. 2025-01-02T15:04:05Z)"})
		return
	}

	notif, err := h.notifier.Create(usecases.CreateRequest{
		Title:       req.Title,
		Message:     req.Message,
		Channel:     domain.Channel(req.Channel),
		Recipient:   req.Recipient,
		ScheduledAt: scheduledAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notif)
}

func (h *Handler) getByID(c *ginext.Context) {
	id := c.Param("id")
	notif, err := h.notifier.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ginext.H{"error": "notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notif)
}

func (h *Handler) cancel(c *ginext.Context) {
	id := c.Param("id")
	err := h.notifier.Cancel(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, ginext.H{"error": "notification not found"})
		case errors.Is(err, domain.ErrAlreadySent):
			c.JSON(http.StatusConflict, ginext.H{"error": "notification already sent"})
		case errors.Is(err, domain.ErrAlreadyCancelled):
			c.JSON(http.StatusConflict, ginext.H{"error": "notification already cancelled"})
		default:
			c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, ginext.H{"message": "notification cancelled"})
}

func (h *Handler) listAll(c *ginext.Context) {
	notifications, err := h.notifier.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}
	if notifications == nil {
		notifications = []*domain.Notification{}
	}
	c.JSON(http.StatusOK, notifications)
}
