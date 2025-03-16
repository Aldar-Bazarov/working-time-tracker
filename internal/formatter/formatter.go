package formatter

import (
	"fmt"
	"time"
	
	"working-time-tracker/internal/config"
)

type Formatter struct {
	config *config.Config
}

func New(config *config.Config) *Formatter {
	return &Formatter{config: config}
}

func (f *Formatter) FormatTime(seconds int64) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	
	if f.config.Display.ShowSeconds {
		return fmt.Sprintf(f.config.Display.TimeFormat, h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", h, m)
}

func (f *Formatter) Today() string {
	return time.Now().Format(f.config.Display.DateFormat)
} 