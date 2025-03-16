package service

import (
	"fmt"
	"time"
	
	"working-time-tracker/internal/models"
	"working-time-tracker/internal/storage"
	"working-time-tracker/internal/formatter"
	"working-time-tracker/internal/logger"
)

type Service struct {
	storage *storage.Storage
	logger  *logger.Logger
}

func New(storage *storage.Storage, logger *logger.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Start() error {
	session, err := s.storage.Load()
	if err != nil {
		return err
	}
	
	today := formatter.Today()
	// ... логика start ...
	return nil
}

func (s *Service) Stop() error {
	// ... логика stop ...
}

func (s *Service) Status() error {
	// ... логика status ...
}

func (s *Service) CleanupStaleSessions(session *models.WorkSession) {
	// ... логика очистки ...
} 