package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"

	"go-vless-client/internal/config"
)

const connectTimeout = 30 * time.Second

// Client управляет жизненным циклом sing-box и хранит состояние подключения.
type Client struct {
	mu     sync.Mutex
	box    *box.Box
	cancel context.CancelFunc
	stats  *Stats
	logger *AppLogger
}

func NewClient(logger *AppLogger) *Client {
	return &Client{
		stats:  NewStats(),
		logger: logger,
	}
}

// Connect запускает sing-box с конфигурацией сервера и настройками режима.
// Возвращает ошибку если запуск не завершился за connectTimeout.
func (c *Client) Connect(ctx context.Context, srv config.ServerConfig, settings config.AppSettings) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.box != nil {
		return fmt.Errorf("already connected")
	}

	c.logger.Add("info", fmt.Sprintf("подключение к %s (%s:%d)…", srv.Name, srv.Address, srv.Port))

	opts, err := buildOptions(srv, settings)
	if err != nil {
		c.logger.Add("error", fmt.Sprintf("ошибка конфигурации: %v", err))
		return fmt.Errorf("build options: %w", err)
	}

	boxCtx, cancel := context.WithCancel(ctx)
	boxCtx = include.Context(boxCtx)

	b, err := box.New(box.Options{
		Context: boxCtx,
		Options: opts,
	})
	if err != nil {
		cancel()
		c.logger.Add("error", fmt.Sprintf("ошибка инициализации: %v", err))
		return fmt.Errorf("create box: %w", err)
	}

	c.logger.Add("info", "запуск sing-box…")

	startErr := make(chan error, 1)
	go func() { startErr <- b.Start() }()

	select {
	case err := <-startErr:
		if err != nil {
			cancel()
			c.logger.Add("error", fmt.Sprintf("ошибка запуска: %v", err))
			return fmt.Errorf("start box: %w", err)
		}
	case <-time.After(connectTimeout):
		cancel()
		b.Close() //nolint:errcheck
		c.logger.Add("error", "превышено время ожидания подключения (30с)")
		return fmt.Errorf("connection timed out after 30s")
	}

	c.box = b
	c.cancel = cancel
	c.stats.Reset()
	c.logger.Add("info", fmt.Sprintf("подключено к %s", srv.Name))
	return nil
}

// Disconnect останавливает sing-box.
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.box == nil {
		return nil
	}

	c.logger.Add("info", "отключение…")
	c.cancel()
	err := c.box.Close()
	c.box = nil
	c.cancel = nil
	if err != nil {
		c.logger.Add("error", fmt.Sprintf("ошибка отключения: %v", err))
	} else {
		c.logger.Add("info", "отключено")
	}
	return err
}

// IsConnected возвращает true если соединение активно.
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.box != nil
}

// Stats возвращает текущие счётчики трафика.
func (c *Client) Stats() config.Stats {
	return c.stats.Get()
}

// Ping измеряет TCP-задержку до сервера. Не требует активного подключения.
func (c *Client) Ping(ctx context.Context, srv config.ServerConfig) (time.Duration, error) {
	addr := fmt.Sprintf("%s:%d", srv.Address, srv.Port)

	start := time.Now()
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("ping %s: %w", addr, err)
	}
	conn.Close()

	return time.Since(start), nil
}

// CheckProxy делает HTTP-запрос через локальный HTTP-прокси и возвращает внешний IP.
// Используется для проверки работоспособности туннеля после подключения.
func (c *Client) CheckProxy(ctx context.Context, settings config.AppSettings) (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected")
	}

	proxyURL := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%d", settings.HTTPPort),
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("check proxy: %w", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 64)
	n, err := resp.Body.Read(buf)
	if err != nil && n == 0 {
		return "", fmt.Errorf("read response: %w", err)
	}

	return string(buf[:n]), nil
}