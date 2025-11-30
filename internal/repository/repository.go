package repository

import "link-service/internal/domain"

type Repository interface {
	Init(dirPath string, fileName string, tempFileName string) error
	SaveRecord(record *domain.Record) error
	SaveTempRecord(record *domain.Record) error
	LoadTempRecords() ([]domain.Record, error)
	ClearTempFile() error
	LoadLastLinksNum() int64
}
