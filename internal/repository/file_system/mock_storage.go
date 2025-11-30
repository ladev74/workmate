package filesystem

import "link-service/internal/domain"

type MockStorage struct{}

func NewMockStorage() *MockStorage { return &MockStorage{} }

func (ms *MockStorage) Init(dirPath string, fileName string, tempFileName string) error { return nil }
func (ms *MockStorage) SaveRecord(record *domain.Record) error                          { return nil }
func (ms *MockStorage) SaveTempRecord(record *domain.Record) error                      { return nil }
func (ms *MockStorage) LoadLastLinksNum() int64                                         { return 0 }
