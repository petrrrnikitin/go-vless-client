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

// Client управляет жизненным циклом sing-box и хранит состояние подключения.
type Client struct {
	mu     sync.Mutex
	box    *box.Box
	cancel context.CancelFunc
	stats  *Stats
}

func NewClient() *Client {
	return &Client{
		stats: NewStats(),
	}
}

// Connect запускает sing-box с конфигурацией сервера и настройками режима.
func (c *Client) Connect(ctx context.Context, srv config.ServerConfig, settings config.AppSettings) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.box != nil {
		return fmt.Errorf("already connected")
	}

	opts, err := buildOptions(srv, settings)
	if err != nil {
		return fmt.Errorf("build options: %w", err)
	}

	boxCtx, cancel := context.WithCancel(ctx)
	boxCtx = include.Context(boxCtx) // регистрируем все протоколы sing-box

	b, err := box.New(box.Options{
		Context: boxCtx,
		Options: opts,
	})
	if err != nil {
		cancel()
		return fmt.Errorf("create box: %w", err)
	}

	if err := b.Start(); err != nil {
		cancel()
		return fmt.Errorf("start box: %w", err)
	}

	c.box = b
	c.cancel = cancel
	c.stats.Reset()
	return nil
}

// Disconnect останавливает sing-box.
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.box == nil {
		return nil
	}

	c.cancel()
	err := c.box.Close()
	c.box = nil
	c.cancel = nil
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