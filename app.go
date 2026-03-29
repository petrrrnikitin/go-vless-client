package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"go-vless-client/internal/api"
	"go-vless-client/internal/config"
	"go-vless-client/internal/core"
)

// App — центральный объект приложения, методы которого доступны из Vue через Wails.
type App struct {
	ctx     context.Context
	storage *config.Storage
	client  *core.Client

	// активное подключение
	activeServerID string
}

func NewApp() *App {
	return &App{
		client: core.NewClient(),
	}
}

// startup вызывается Wails при старте приложения.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	storage, err := config.NewStorage()
	if err != nil {
		runtime.LogErrorf(ctx, "failed to init storage: %v", err)
		return
	}
	a.storage = storage

	// запускаем REST API для браузерного расширения
	apiSrv := api.NewServer(a)
	if err := apiSrv.Start(ctx, storage.Settings().APIPort); err != nil {
		runtime.LogErrorf(ctx, "failed to start api server: %v", err)
	}

	// запускаем фоновую отправку статистики
	go a.statsLoop(ctx)

	// авто-подключение к последнему серверу
	settings := storage.Settings()
	if settings.AutoConnect && settings.LastServerID != "" {
		go func() {
			if err := a.Connect(settings.LastServerID); err != nil {
				runtime.LogErrorf(ctx, "auto-connect failed: %v", err)
			}
		}()
	}
}

// shutdown вызывается Wails при закрытии приложения.
func (a *App) shutdown(_ context.Context) {
	if a.client.IsConnected() {
		a.client.Disconnect() //nolint:errcheck
	}
}

// --- Серверы ---

func (a *App) GetServers() []config.ServerConfig {
	return a.storage.Servers()
}

func (a *App) SaveServer(cfg config.ServerConfig) error {
	if cfg.ID == "" {
		cfg.ID = uuid.NewString()
	}
	return a.storage.SaveServer(cfg)
}

func (a *App) DeleteServer(id string) error {
	return a.storage.DeleteServer(id)
}

// --- Подключение ---

func (a *App) Connect(serverID string) error {
	srv, ok := a.storage.ServerByID(serverID)
	if !ok {
		return fmt.Errorf("server not found: %s", serverID)
	}

	settings := a.storage.Settings()

	if err := a.client.Connect(a.ctx, srv, settings); err != nil {
		return err
	}

	a.activeServerID = serverID

	// сохраняем последний использованный сервер
	settings.LastServerID = serverID
	a.storage.SaveSettings(settings) //nolint:errcheck

	a.emitStatus()
	return nil
}

func (a *App) Disconnect() error {
	err := a.client.Disconnect()
	a.activeServerID = ""
	a.emitStatus()
	return err
}

func (a *App) GetStatus() config.ConnectionStatus {
	if !a.client.IsConnected() {
		return config.ConnectionStatus{
			Connected: false,
			Mode:      a.storage.Settings().Mode,
		}
	}

	srv, _ := a.storage.ServerByID(a.activeServerID)
	return config.ConnectionStatus{
		Connected:  true,
		ServerID:   a.activeServerID,
		ServerName: srv.Name,
		Mode:       a.storage.Settings().Mode,
	}
}

// --- Настройки ---

func (a *App) GetSettings() config.AppSettings {
	return a.storage.Settings()
}

func (a *App) SaveSettings(s config.AppSettings) error {
	return a.storage.SaveSettings(s)
}

// --- Импорт / Экспорт ---

// ParseURI разбирает VLESS URI и возвращает конфигурацию сервера для предпросмотра.
// ID не заполняется — сервер не сохраняется, пользователь должен явно вызвать SaveServer.
func (a *App) ParseURI(uri string) (config.ServerConfig, error) {
	return core.ParseVLESSURI(uri)
}

// ExportURI возвращает VLESS URI для указанного сервера.
func (a *App) ExportURI(serverID string) (string, error) {
	srv, ok := a.storage.ServerByID(serverID)
	if !ok {
		return "", fmt.Errorf("server not found: %s", serverID)
	}
	return core.BuildVLESSURI(srv), nil
}

// --- Проверка соединения ---

// Ping измеряет TCP-задержку до сервера в миллисекундах.
func (a *App) Ping(serverID string) (int64, error) {
	srv, ok := a.storage.ServerByID(serverID)
	if !ok {
		return 0, fmt.Errorf("server not found: %s", serverID)
	}

	d, err := a.client.Ping(a.ctx, srv)
	if err != nil {
		return 0, err
	}
	return d.Milliseconds(), nil
}

// CheckProxy проверяет работу туннеля и возвращает внешний IP.
func (a *App) CheckProxy() (string, error) {
	return a.client.CheckProxy(a.ctx, a.storage.Settings())
}

// --- Events ---

func (a *App) emitStatus() {
	runtime.EventsEmit(a.ctx, "status:changed", a.GetStatus())
}

// statsLoop каждые 2 секунды отправляет статистику трафика во фронтенд.
func (a *App) statsLoop(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if a.client.IsConnected() {
				runtime.EventsEmit(ctx, "stats:update", a.client.Stats())
			}
		}
	}
}
