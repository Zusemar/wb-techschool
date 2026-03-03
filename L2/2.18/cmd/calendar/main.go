package main

import (
	"log"
	"net/http"

	"calendar/internal/handlers"
	"calendar/internal/repo"
	"calendar/internal/usecases"
)

func main() {
	// repository layer
	repository := repo.NewMemoryRepository()

	// usecase layer
	calendarUsecase := usecases.NewCalendar(repository)

	// handler layer
	handler := handlers.NewHandler(calendarUsecase)

	// router
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", handler.CreateEvent)
	mux.HandleFunc("/update_event", handler.UpdateEvent)
	mux.HandleFunc("/delete_event", handler.DeleteEvent)
	mux.HandleFunc("/events_for_day", handler.EventsForDay)
	mux.HandleFunc("/events_for_week", handler.EventsForWeek)
	mux.HandleFunc("/events_for_month", handler.EventsForMonth)

	// middleware wrapping
	wrapped := handlers.RecoveryMiddleware(
		handlers.LoggingMiddleware(mux),
	)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", wrapped); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
