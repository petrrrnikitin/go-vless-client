package core

import (
	"context"
	"fmt"
	"sync"

	box "github.com/sagernet/sing-box"

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