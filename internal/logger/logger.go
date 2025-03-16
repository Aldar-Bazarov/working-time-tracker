package logger

import "fmt"

type Logger struct {
    enabled bool
}

func New(enabled bool) *Logger {
    return &Logger{enabled: enabled}
}

func (l *Logger) Log(format string, args ...interface{}) {
    if l.enabled {
        fmt.Printf("[DEBUG] "+format+"\n", args...)
    }
} 