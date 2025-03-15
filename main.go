package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var dataFile = filepath.Join(os.Getenv("HOME"), ".wtt_data.json")

type WorkSession struct {
	StartTime  time.Time `json:"start_time"`
	TotalTime  int64     `json:"total_time"`
	IsRunning  bool      `json:"is_running"`
	IsPaused   bool      `json:"is_paused"`
	LastUpdate time.Time `json:"last_update"`
}

func loadData() WorkSession {
	var session WorkSession
	file, err := os.ReadFile(dataFile)
	if err == nil {
		json.Unmarshal(file, &session)
	}
	return session
}

func saveData(session WorkSession) {
	data, _ := json.MarshalIndent(session, "", "  ")
	os.WriteFile(dataFile, data, 0o644)
}

func start() {
	session := loadData()
	if session.IsRunning {
		fmt.Println("Работа уже запущена!")
		return
	}
	session.StartTime = time.Now()
	session.IsRunning = true
	session.IsPaused = false
	session.LastUpdate = time.Now()
	saveData(session)
	fmt.Println("Работа начата!")
}

func pause() {
	session := loadData()
	if !session.IsRunning || session.IsPaused {
		fmt.Println("Нельзя поставить на паузу!")
		return
	}
	session.TotalTime += int64(time.Since(session.LastUpdate).Seconds())
	session.IsPaused = true
	session.LastUpdate = time.Now()
	saveData(session)
	fmt.Println("Работа на паузе!")
}

func resume() {
	session := loadData()
	if !session.IsPaused {
		fmt.Println("Работа не на паузе!")
		return
	}
	session.IsPaused = false
	session.LastUpdate = time.Now()
	saveData(session)
	fmt.Println("Работа продолжается!")
}

func stop() {
	session := loadData()
	if !session.IsRunning {
		fmt.Println("Нет активной работы!")
		return
	}
	session.TotalTime += int64(time.Since(session.LastUpdate).Seconds())
	session.IsRunning = false
	session.IsPaused = false
	saveData(session)
	fmt.Printf("Работа завершена! Всего отработано: %s\n", formatTime(session.TotalTime))
}

func show() {
	session := loadData()
	status := "Остановлено"
	if session.IsRunning {
		status = "Работает"
	} else if session.IsPaused {
		status = "На паузе"
	}
	fmt.Printf("Статус: %s\n", status)
	fmt.Printf("Общее время: %s\n", formatTime(session.TotalTime))
}

func daemonize() {
	cmd := exec.Command(os.Args[0], "daemon")
	cmd.Start()
	fmt.Println("Демон запущен!")
}

func handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Получен сигнал завершения, сохраняю данные...")
		stop()
		os.Exit(0)
	}()
}

func installSystemdService() {
	serviceContent := `[Unit]
Description=Working Time Tracker
After=network.target

[Service]
ExecStart=` + os.Args[0] + ` daemon
Restart=always
User=` + os.Getenv("USER") + `

[Install]
WantedBy=multi-user.target`

	servicePath := "/etc/systemd/system/wtt.service"
	err := os.WriteFile(servicePath, []byte(serviceContent), 0644)
	if err != nil {
		fmt.Println("Ошибка установки systemd сервиса:", err)
		return
	}

	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", "wtt").Run()
	exec.Command("systemctl", "start", "wtt").Run()

	fmt.Println("Systemd сервис установлен и запущен!")
}

func formatTime(seconds int64) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println(os.Args)
		fmt.Println("Использование: wtt <start|pause|resume|stop|show|daemon|install-service>")
		return
	}

	handleSignals()

	switch os.Args[1] {
	case "start":
		start()
	case "pause":
		pause()
	case "resume":
		resume()
	case "stop":
		stop()
	case "show":
		show()
	case "daemon":
		daemonize()
	case "install-service":
		installSystemdService()
	default:
		fmt.Println("Неизвестная команда")
	}
}
