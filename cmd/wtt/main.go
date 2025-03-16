package main

import (
	"fmt"
	"os"
	
	"working-time-tracker/internal/config"
	"working-time-tracker/internal/storage"
	"working-time-tracker/internal/logger"
	"working-time-tracker/internal/service"
	"working-time-tracker/internal/formatter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Использование: wtt <start|stop|status>")
		return
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	store := storage.New(cfg)
	log := logger.New(cfg.Logging.Enabled)
	fmt := formatter.New(cfg)
	svc := service.New(store, log, fmt)

	var err error
	switch os.Args[1] {
	case "start":
		err = svc.Start()
	case "stop":
		err = svc.Stop()
	case "status":
		err = svc.Status()
	default:
		fmt.Println("Неизвестная команда")
		return
	}

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}
} 