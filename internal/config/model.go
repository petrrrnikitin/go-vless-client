package config

// Mode определяет режим работы приложения.
type Mode string

const (
	ModeProxy Mode = "proxy" // только SOCKS5/HTTP прокси
	ModeVPN   Mode = "vpn"   // только TUN (весь трафик системы)
	ModeBoth  Mode = "both"  // прокси + TUN одновременно
)

// Transport определяет транспортный протокол VLESS.
type Transport string

const (
	TransportTCP Transport = "tcp"
	TransportWS  Transport = "ws"
)

// ServerConfig хранит параметры одного VLESS-сервера.
type ServerConfig struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Port      int       `json:"port"`
	UUID      string    `json:"uuid"`
	Transport Transport `json:"transport"`
	TLS       bool      `json:"tls"`
	SNI       string    `json:"sni,omitempty"`  // имя сервера для TLS
	Path      string    `json:"path,omitempty"` // путь для WebSocket
	Flow      string    `json:"flow,omitempty"` // VLESS flow, например "xtls-rprx-vision"
}

// AppSettings хранит настройки приложения.
type AppSettings struct {
	Mode         Mode   `json:"mode"`
	Socks5Port   int    `json:"socks5_port"`
	HTTPPort     int    `json:"http_port"`
	APIPort      int    `json:"api_port"`
	LastServerID string `json:"last_server_id,omitempty"`
	AutoConnect  bool   `json:"auto_connect"`
}

// DefaultSettings возвращает настройки по умолчанию.
func DefaultSettings() AppSettings {
	return AppSettings{
		Mode:       ModeProxy,
		Socks5Port: 1080,
		HTTPPort:   8080,
		APIPort:    9090,
	}
}

// ConnectionStatus описывает текущее состояние подключения.
type ConnectionStatus struct {
	Connected  bool   `json:"connected"`
	ServerID   string `json:"server_id,omitempty"`
	ServerName string `json:"server_name,omitempty"`
	Mode       Mode   `json:"mode"`
}

// Stats содержит счётчики трафика текущей сессии.
type Stats struct {
	Upload   int64 `json:"upload"`   // байт отправлено
	Download int64 `json:"download"` // байт получено
}
