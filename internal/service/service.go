package service

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"link-service/internal/domain"
	"link-service/internal/repository"
)

const (
	StatusAvailable    = "available"
	StatusNotAvailable = "not available"
)

type Service struct {
	repository  repository.Repository
	mu          *sync.Mutex
	counter     int
	pingTimeout time.Duration
	logger      *zap.Logger
}

// TODO: counter logic

func New(repo repository.Repository, pingTimeout time.Duration, logger *zap.Logger) *Service {
	return &Service{
		repository:  repo,
		mu:          &sync.Mutex{},
		counter:     0,
		pingTimeout: pingTimeout,
		logger:      logger,
	}
}

func (s *Service) Process(links []string) error {
	client := http.Client{
		Timeout: s.pingTimeout,
	}

	rec := &domain.Record{
		Links: make(map[string]string),
		ID:    0,
	}

	for _, link := range links {
		// TODO: const
		resp, err := client.Head("https://" + link)
		if err != nil {
			s.logger.Error("failed to ping link", zap.String("link", link), zap.Error(err))
			return fmt.Errorf("failed to ping link: %s: %w", link, err)
		}

		if resp.StatusCode == http.StatusOK {
			rec.Links[link] = StatusAvailable
		} else {
			rec.Links[link] = StatusNotAvailable
		}
	}

	s.incCounter()
	rec.ID = s.counter

	err := s.repository.SaveRecord(rec)
	if err != nil {
		s.decCounter()

		s.logger.Error("failed to save record", zap.Error(err))
		return fmt.Errorf("failed to save record: %w", err)
	}

	fmt.Println(rec)
	return nil
}

func (s *Service) incCounter() {
	s.mu.Lock()
	s.counter++
	s.mu.Unlock()
}

func (s *Service) decCounter() {
	s.mu.Lock()
	s.counter--
	s.mu.Unlock()
}
