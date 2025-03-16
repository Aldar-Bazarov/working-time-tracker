package config

import (
    "fmt"
    "os"
    "path/filepath"
    
    "gopkg.in/yaml.v3"
)

type Config struct {
    Storage StorageConfig `yaml:"storage"`
    Logging LoggingConfig `yaml:"logging"`
    Display DisplayConfig `yaml:"display"`
    Backup  BackupConfig  `yaml:"backup"`
}

type StorageConfig struct {
    DataFilePath string `yaml:"data_file_path"`
}

type LoggingConfig struct {
    Enabled bool   `yaml:"enabled"`
    Level   string `yaml:"level"`
}

type DisplayConfig struct {
    TimeFormat     string `yaml:"time_format"`
    DateFormat     string `yaml:"date_format"`
    ShowSeconds    bool   `yaml:"show_seconds"`
    CompactOutput  bool   `yaml:"compact_output"`
}

type BackupConfig struct {
    Enabled        bool   `yaml:"enabled"`
    IntervalHours  int    `yaml:"interval_hours"`
    MaxBackups     int    `yaml:"max_backups"`
    BackupDir      string `yaml:"backup_dir"`
}

func getDefaultConfig() Config {
    home, _ := os.UserHomeDir()
    
    return Config{
        Storage: StorageConfig{
            DataFilePath: filepath.Join(home, ".wtt_data.json"),
        },
        Logging: LoggingConfig{
            Enabled: false,
            Level:   "info",
        },
        Display: DisplayConfig{
            TimeFormat:    "%02d:%02d:%02d",
            DateFormat:    "2006-01-02",
            ShowSeconds:   true,
            CompactOutput: false,
        },
        Backup: BackupConfig{
            Enabled:       true,
            IntervalHours: 24,
            MaxBackups:    7,
            BackupDir:     filepath.Join(home, ".wtt_backups"),
        },
    }
}

func Load() (*Config, error) {
    config := getDefaultConfig()
    
    configPath := os.Getenv("WTT_CONFIG")
    if configPath == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            return nil, fmt.Errorf("не удалось получить домашнюю директорию: %w", err)
        }
        configPath = filepath.Join(home, ".wtt_config.yaml")
    }
    
    // Если файл конфигурации не существует, создаем его с дефолтными настройками
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        if err := config.Save(configPath); err != nil {
            return nil, fmt.Errorf("не удалось создать файл конфигурации: %w", err)
        }
        return &config, nil
    }
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
    }
    
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
    }
    
    return &config, nil
}

func (c *Config) Save(path string) error {
    data, err := yaml.Marshal(c)
    if err != nil {
        return fmt.Errorf("ошибка сериализации конфигурации: %w", err)
    }
    
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        return fmt.Errorf("не удалось создать директорию конфигурации: %w", err)
    }
    
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("не удалось сохранить файл конфигурации: %w", err)
    }
    
    return nil
} 