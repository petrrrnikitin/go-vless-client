package core

import (
	"fmt"
	"net/url"
	"strconv"

	"go-vless-client/internal/config"
)

// BuildVLESSURI собирает vless:// URI из конфигурации сервера.
func BuildVLESSURI(cfg config.ServerConfig) string {
	q := url.Values{}

	if cfg.TLS {
		q.Set("security", "tls")
	}
	if cfg.SNI != "" {
		q.Set("sni", cfg.SNI)
	}
	if cfg.Transport == config.TransportWS {
		q.Set("type", "ws")
		if cfg.Path != "" && cfg.Path != "/" {
			q.Set("path", cfg.Path)
		}
	} else {
		q.Set("type", "tcp")
	}
	if cfg.Flow != "" {
		q.Set("flow", cfg.Flow)
	}

	u := url.URL{
		Scheme:   "vless",
		User:     url.User(cfg.UUID),
		Host:     fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		RawQuery: q.Encode(),
		Fragment: cfg.Name,
	}
	return u.String()
}

// ParseVLESSURI разбирает ссылку вида vless://uuid@host:port?params#name
// в конфигурацию сервера. Поддерживаются транспорты tcp и ws, безопасность tls.
func ParseVLESSURI(uri string) (config.ServerConfig, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return config.ServerConfig{}, fmt.Errorf("parse uri: %w", err)
	}
	if u.Scheme != "vless" {
		return config.ServerConfig{}, fmt.Errorf("неподдерживаемая схема: %s", u.Scheme)
	}

	uuid := u.User.Username()
	if uuid == "" {
		return config.ServerConfig{}, fmt.Errorf("пустой UUID")
	}

	host := u.Hostname()
	if host == "" {
		return config.ServerConfig{}, fmt.Errorf("пустой адрес сервера")
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil || port <= 0 || port > 65535 {
		return config.ServerConfig{}, fmt.Errorf("неверный порт: %s", u.Port())
	}

	q := u.Query()

	transport := config.TransportTCP
	if q.Get("type") == "ws" {
		transport = config.TransportWS
	}

	name := u.Fragment
	if name == "" {
		name = host
	}

	path := q.Get("path")
	if transport == config.TransportWS && path == "" {
		path = "/"
	}

	return config.ServerConfig{
		Name:      name,
		Address:   host,
		Port:      port,
		UUID:      uuid,
		Transport: transport,
		TLS:       q.Get("security") == "tls",
		SNI:       q.Get("sni"),
		Path:      path,
		Flow:      q.Get("flow"),
	}, nil
}
