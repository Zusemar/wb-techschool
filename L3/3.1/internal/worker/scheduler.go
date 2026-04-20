package worker

import (
	"context"
	"log"
	"time"

	"notifier/internal/domain"
	"notifier/internal/queue"
	"notifier/internal/repo"
)

// Scheduler polls for pending notifications and publishes them to RabbitMQ
// when their scheduled time has been reached.
type Scheduler struct {
	repo      repo.Repository
	publisher *queue.Publisher
	interval  time.Duration
}

func NewScheduler(r repo.Repository, p *queue.Publisher) *Scheduler {
	return &Scheduler{
		repo:      r,
		publisher: p,
		interval:  time.Second,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.dispatch(ctx)
		}
	}
}

func (s *Scheduler) dispatch(ctx context.Context) {
	pending, err := s.repo.ListPending()
	if err != nil {
		log.Printf("ERROR: scheduler list pending: %v", err)
		return
	}

	now := time.Now()
	for _, n := range pending {
		if !n.ScheduledAt.After(now) {
			if err := s.publisher.Publish(ctx, n); err != nil {
				log.Printf("ERROR: scheduler publish %s: %v", n.ID, err)
				continue
			}
			if err := s.repo.UpdateStatus(n.ID, domain.StatusQueued); err != nil {
				log.Printf("ERROR: scheduler update status %s: %v", n.ID, err)
			} else {
				log.Printf("INFO: notification %s queued for delivery", n.ID)
			}
		}
	}
}
