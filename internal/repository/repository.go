package repository

import "link-service/internal/domain"

type Repository interface {
	Init(dirPath string, fileName string, tempFileName string) error
	SaveRecord(record *domain.Record) error
	SaveTempRecord(record *domain.Record) error
	LoadLastLinksNum() int64
}
