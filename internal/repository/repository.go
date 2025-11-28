package repository

import "link-service/internal/domain"

type Repository interface {
	Init(path string, fileName string) error
	SaveRecord(record *domain.Record) error
}
