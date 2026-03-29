package core

import (
	"fmt"
	"net/url"
	"strconv"

	"go-vless-client/internal/config"
)

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
