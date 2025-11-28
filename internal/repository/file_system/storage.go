package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"link-service/internal/domain"
)

type Storage struct {
	mu     *sync.Mutex
	path   string
	logger *zap.Logger
}

func New(logger *zap.Logger) *Storage {
	return &Storage{
		mu:     &sync.Mutex{},
		logger: logger,
	}
}

func (s *Storage) Init(path string, fileName string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		s.logger.Error("failed to create dir", zap.String("path", path), zap.Error(err))
		return fmt.Errorf("failed to create dir: %s: %w", path, err)
	}

	file, err := os.OpenFile(filepath.Join(path, fileName),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		s.logger.Error("failed to create file", zap.String("file_name", fileName), zap.Error(err))
		return fmt.Errorf("failed to create file: %s: %w", fileName, err)
	}

	defer file.Close()

	s.path = filepath.Join(path, fileName)

	s.logger.Info("file created", zap.String("file", fileName))
	return nil
}

func (s *Storage) SaveRecord(record *domain.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		s.logger.Error("failed to open file", zap.String("path", s.path), zap.Error(err))
		return fmt.Errorf("failed to open file: %s: %w", s.path, err)
	}
	defer file.Close()

	data, err := json.Marshal(&record)
	if err != nil {
		s.logger.Error("failed to marshal record", zap.Error(err))
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	data = append(data, '\n')

	_, err = file.Write(data)
	if err != nil {
		s.logger.Error("failed to write record", zap.Error(err))
		return fmt.Errorf("failed to write record: %w", err)
	}

	s.logger.Info("successfully wrote record")
	return nil
}
