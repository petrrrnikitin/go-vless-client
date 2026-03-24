package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"go-vless-client/internal/config"
)

// StatusResponse — ответ на GET /status.
type StatusResponse struct {
	Connected   bool   `json:"connected"`
	Mode        string `json:"mode"`
	Server      string `json:"server,omitempty"`
	Socks5Proxy string `json:"socks5_proxy"`
	HTTPProxy   string `json:"http_proxy"`
}

// StatusProvider описывает зависимость сервера от App.
type StatusProvider interface {
	GetStatus() config.ConnectionStatus
	GetSettings() config.AppSettings
}

// Server — локальный HTTP-сервер для браузерного расширения.
// Слушает только на loopback, не доступен извне.
type Server struct {
	provider StatusProvider
	srv      *http.Server
}

func NewServer(provider StatusProvider) *Server {
	s := &Server{provider: provider}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", s.handleStatus)

	s.srv = &http.Server{Handler: mux}
	return s
}

// Start запускает сервер на указанном порту.
func (s *Server) Start(ctx context.Context, port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("listen api: %w", err)
	}

	go func() {
		<-ctx.Done()
		s.srv.Close() //nolint:errcheck
	}()

	go func() {
		s.srv.Serve(ln) //nolint:errcheck
	}()

	return nil
}

func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	status := s.provider.GetStatus()
	settings := s.provider.GetSettings()

	resp := StatusResponse{
		Connected:   status.Connected,
		Mode:        string(status.Mode),
		Server:      status.ServerName,
		Socks5Proxy: fmt.Sprintf("127.0.0.1:%d", settings.Socks5Port),
		HTTPProxy:   fmt.Sprintf("127.0.0.1:%d", settings.HTTPPort),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}