package models

import "time"

type WorkDay struct {
    Date       string    `json:"date"`      
    TotalTime  int64     `json:"total_time"`
    IsRunning  bool      `json:"is_running"`
    LastUpdate time.Time `json:"last_update"`
}

type WorkSession struct {
    Days       []WorkDay `json:"days"`    
    CurrentDay *WorkDay  `json:"current_day"`
}

func (w *WorkDay) Validate() error {
    if w.TotalTime < 0 {
        return fmt.Errorf("отрицательное время работы")
    }
    if w.IsRunning && w.LastUpdate.IsZero() {
        return fmt.Errorf("активная сессия без времени последнего обновления")
    }
    return nil
} 