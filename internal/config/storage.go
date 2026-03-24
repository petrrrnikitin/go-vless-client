package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const appDirName = "go-vless-client"

// configData — структура JSON-файла на диске.
type configData struct {
	Servers  []ServerConfig `json:"servers"`
	Settings AppSettings    `json:"settings"`
}

// Storage управляет хранением конфигурации на диске.
type Storage struct {
	path string
	data configData
}

// NewStorage создаёт Storage и загружает конфиг с диска.
// Если файл не существует — создаёт с настройками по умолчанию.
func NewStorage() (*Storage, error) {
	dir, err := configDir()
	if err != nil {
		return nil, fmt.Errorf("config dir: %w", err)
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	s := &Storage{
		path: filepath.Join(dir, "config.json"),
		data: configData{
			Settings: DefaultSettings(),
		},
	}

	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return s, nil
}

// Servers возвращает список серверов.
func (s *Storage) Servers() []ServerConfig {
	if s.data.Servers == nil {
		return []ServerConfig{}
	}
	return s.data.Servers
}

// ServerByID возвращает сервер по ID.
func (s *Storage) ServerByID(id string) (ServerConfig, bool) {
	for _, srv := range s.data.Servers {
		if srv.ID == id {
			return srv, true
		}
	}
	return ServerConfig{}, false
}

// SaveServer добавляет новый или обновляет существующий сервер.
func (s *Storage) SaveServer(cfg ServerConfig) error {
	for i, srv := range s.data.Servers {
		if srv.ID == cfg.ID {
			s.data.Servers[i] = cfg
			return s.save()
		}
	}
	s.data.Servers = append(s.data.Servers, cfg)
	return s.save()
}

// DeleteServer удаляет сервер по ID.
func (s *Storage) DeleteServer(id string) error {
	servers := make([]ServerConfig, 0, len(s.data.Servers))
	for _, srv := range s.data.Servers {
		if srv.ID != id {
			servers = append(servers, srv)
		}
	}
	s.data.Servers = servers
	return s.save()
}

// Settings возвращает текущие настройки приложения.
func (s *Storage) Settings() AppSettings {
	return s.data.Settings
}

// SaveSettings сохраняет настройки приложения.
func (s *Storage) SaveSettings(settings AppSettings) error {
	s.data.Settings = settings
	return s.save()
}

func (s *Storage) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&s.data)
}

func (s *Storage) save() error {
	f, err := os.CreateTemp(filepath.Dir(s.path), "config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := f.Name()

	if err := json.NewEncoder(f).Encode(s.data); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("encode config: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, s.path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("replace config file: %w", err)
	}

	return nil
}

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appDirName), nil
}