package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"working-time-tracker/internal/config"
	"working-time-tracker/internal/models"
)

type Storage struct {
	config *config.Config
}

func New(config *config.Config) *Storage {
	return &Storage{config: config}
}

func (s *Storage) Load() (models.WorkSession, error) {
	var session models.WorkSession
	file, err := os.ReadFile(s.config.Storage.DataFilePath)
	if err != nil && !os.IsNotExist(err) {
		return session, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	
	if len(file) > 0 {
		if err := json.Unmarshal(file, &session); err != nil {
			return session, fmt.Errorf("ошибка десериализации: %w", err)
		}
	}
	
	if session.Days == nil {
		session.Days = make([]models.WorkDay, 0)
	}
	
	return session, nil
}

func (s *Storage) Save(session models.WorkSession) error {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}
	
	if err := os.WriteFile(s.config.Storage.DataFilePath, data, 0o644); err != nil {
		return fmt.Errorf("ошибка сохранения: %w", err)
	}
	
	return nil
}

func (s *Storage) Backup() error {
	if !s.config.Backup.Enabled {
		return nil
	}
	
	if err := os.MkdirAll(s.config.Backup.BackupDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию для бэкапов: %w", err)
	}
	
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(s.config.Backup.BackupDir, fmt.Sprintf("wtt_backup_%s.json", timestamp))
	
	// Копируем текущий файл данных
	data, err := os.ReadFile(s.config.Storage.DataFilePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла данных: %w", err)
	}
	
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("ошибка создания бэкапа: %w", err)
	}
	
	// Удаляем старые бэкапы
	s.cleanOldBackups()
	
	return nil
}

func (s *Storage) cleanOldBackups() error {
	files, err := os.ReadDir(s.config.Backup.BackupDir)
	if err != nil {
		return err
	}
	
	if len(files) <= s.config.Backup.MaxBackups {
		return nil
	}
	
	// ... логика удаления старых бэкапов ...
	return nil
} 