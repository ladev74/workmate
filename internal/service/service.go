package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"link-service/internal/domain"
	"link-service/internal/repository"
)

const (
	statusAvailable    = "available"
	statusNotAvailable = "not available"
	statusUnknown      = "unknown"

	httpsPrefix = "https://"
	httpPrefix  = "http://"
)

var (
	ErrAppStopped = errors.New("application is stopped")
)

type Service struct {
	counter     int64
	repository  repository.Repository
	httpClient  *http.Client
	pingTimeout time.Duration
	logger      *zap.Logger
}

func New(repo repository.Repository, pingTimeout time.Duration, logger *zap.Logger) *Service {
	lastLinksNum := repo.LoadLastLinksNum()

	return &Service{
		repository:  repo,
		counter:     lastLinksNum,
		httpClient:  &http.Client{Timeout: pingTimeout},
		pingTimeout: pingTimeout,
		logger:      logger,
	}
}

func (s *Service) Process(serverCtx context.Context, requestCtx context.Context, links []string) (*domain.Record, error) {
	s.incCounter()
	rec := &domain.Record{
		Links: make(map[string]string),
		ID:    s.counter,
	}

	select {
	case <-serverCtx.Done():
		for _, link := range links {
			rec.Links[link] = statusUnknown
		}

		err := s.repository.SaveTempRecord(rec)
		if err != nil {
			s.decCounter()
			s.logger.Error("failed to save temp record", zap.Error(err))
			return nil, fmt.Errorf("failed to save temp record: %w", err)
		}

		return rec, ErrAppStopped

	default:
	}

	for _, link := range links {
		select {
		case <-requestCtx.Done():
			s.decCounter()
			s.logger.Info(requestCtx.Err().Error(), zap.String("link", link))
			return nil, requestCtx.Err()

		default:
		}

		statusCode, err := s.ping(link)
		if err != nil || statusCode != http.StatusOK {
			s.logger.Warn("failed to ping link", zap.String("link", link), zap.Error(err))
			rec.Links[link] = statusNotAvailable
		} else {
			rec.Links[link] = statusAvailable
		}
	}

	err := s.repository.SaveRecord(rec)
	if err != nil {
		s.decCounter()
		s.logger.Error("failed to save record", zap.Error(err))
		return nil, fmt.Errorf("failed to save record: %w", err)
	}

	s.logger.Info("success process record")
	return rec, nil
}

func (s *Service) ping(link string) (int, error) {
	if !strings.HasPrefix(link, httpPrefix) && !strings.HasPrefix(link, httpsPrefix) {
		link = httpsPrefix + link
	}

	resp, err := s.httpClient.Head(link)
	if err == nil {
		return resp.StatusCode, nil
	}

	resp, err = s.httpClient.Get(link)
	if err != nil {
		return 0, fmt.Errorf("failed to ping link: %w", err)
	}

	return resp.StatusCode, nil
}

func (s *Service) incCounter() {
	atomic.AddInt64(&s.counter, 1)
}

func (s *Service) decCounter() {
	atomic.AddInt64(&s.counter, -1)
}
