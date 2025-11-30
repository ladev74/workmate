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
	mu       *sync.Mutex
	path     string
	tempPath string
	logger   *zap.Logger
}

func New(logger *zap.Logger) *Storage {
	return &Storage{
		mu:     &sync.Mutex{},
		logger: logger,
	}
}

func (s *Storage) Init(dirPath string, fileName string, tempFileName string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		s.logger.Error("failed to create dir", zap.String("dir_path", dirPath), zap.Error(err))
		return fmt.Errorf("failed to create dir: %s: %w", dirPath, err)
	}

	filePath := filepath.Join(dirPath, fileName)
	tempFilePath := filepath.Join(dirPath, tempFileName)

	file, err := os.OpenFile(filePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		s.logger.Error("failed to create file", zap.String("file_name", fileName), zap.Error(err))
		return fmt.Errorf("failed to create file: %s: %w", fileName, err)
	}

	defer file.Close()

	s.path = filePath

	tempFile, err := os.OpenFile(tempFilePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		s.logger.Error("failed to create temp file", zap.String("file_name", tempFileName), zap.Error(err))
		return fmt.Errorf("failed to create temp file: %s: %w", tempFileName, err)
	}

	defer tempFile.Close()

	s.tempPath = tempFilePath

	s.logger.Info("files created", zap.String("file", filePath), zap.String("temp_path", tempFilePath))
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

	data, err := json.Marshal(record)
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

func (s *Storage) SaveTempRecord(record *domain.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tempFile, err := os.OpenFile(s.tempPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		s.logger.Error("failed to open temp file", zap.String("path", s.tempPath), zap.Error(err))
		return fmt.Errorf("failed to open temp file: %s: %w", s.tempPath, err)
	}
	defer tempFile.Close()

	data, err := json.Marshal(record)
	if err != nil {
		s.logger.Error("failed to marshal temp record", zap.Error(err))
		return fmt.Errorf("failed to marshal temp record: %w", err)
	}

	data = append(data, '\n')

	_, err = tempFile.Write(data)
	if err != nil {
		s.logger.Error("failed to write temp record", zap.Error(err))
		return fmt.Errorf("failed to write temp record: %w", err)
	}

	s.logger.Info("successfully wrote temp record")
	return nil
}

func (s *Storage) LoadLastLinksNum() int64 {
	return 0
}
