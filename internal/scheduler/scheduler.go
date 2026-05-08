package scheduler

import (
	"log"
	"time"

	"github.com/qiblatdigital/zoztool-api/internal/service"
)

type Scheduler struct {
	postSvc  *service.PostService
	interval time.Duration
	stop     chan struct{}
}

func NewScheduler(postSvc *service.PostService) *Scheduler {
	return &Scheduler{
		postSvc:  postSvc,
		interval: 60 * time.Second,
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.run()
			case <-s.stop:
				return
			}
		}
	}()
	log.Println("scheduler started, polling every 60s")
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) run() {
	s.postSvc.PublishScheduled()
}
