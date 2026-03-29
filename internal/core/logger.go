package core

import (
	"sync"
	"time"

	boxlog "github.com/sagernet/sing-box/log"

	"go-vless-client/internal/config"
)

const maxLogEntries = 500

// AppLogger — кольцевой буфер логов. Реализует log.PlatformWriter для sing-box,
// чтобы перехватывать его внутренние сообщения.
type AppLogger struct {
	mu       sync.Mutex
	entries  []config.LogEntry
	onChange func(config.LogEntry)
}

func NewAppLogger() *AppLogger {
	return &AppLogger{}
}

// SetOnChange устанавливает callback, вызываемый при каждой новой записи.
func (l *AppLogger) SetOnChange(fn func(config.LogEntry)) {
	l.mu.Lock()
	l.onChange = fn
	l.mu.Unlock()
}

// Add добавляет запись с заданным уровнем.
func (l *AppLogger) Add(level, message string) {
	entry := config.LogEntry{
		Time:    time.Now().Format("15:04:05"),
		Level:   level,
		Message: message,
	}
	l.mu.Lock()
	l.entries = append(l.entries, entry)
	if len(l.entries) > maxLogEntries {
		l.entries = l.entries[len(l.entries)-maxLogEntries:]
	}
	cb := l.onChange
	l.mu.Unlock()
	if cb != nil {
		cb(entry)
	}
}

// WriteMessage реализует log.PlatformWriter — принимает логи от sing-box.
func (l *AppLogger) WriteMessage(level boxlog.Level, message string) {
	lvl := "info"
	switch level {
	case boxlog.LevelError, boxlog.LevelFatal, boxlog.LevelPanic:
		lvl = "error"
	case boxlog.LevelWarn:
		lvl = "warn"
	case boxlog.LevelDebug, boxlog.LevelTrace:
		lvl = "debug"
	}
	l.Add(lvl, message)
}

// Entries возвращает копию всех накопленных записей.
func (l *AppLogger) Entries() []config.LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	result := make([]config.LogEntry, len(l.entries))
	copy(result, l.entries)
	return result
}
