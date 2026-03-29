package core

import (
	"testing"

	"go-vless-client/internal/config"
)

func TestParseVLESSURI_TCP_TLS(t *testing.T) {
	uri := "vless://550e8400-e29b-41d4-a716-446655440000@example.com:443?security=tls&sni=example.com&flow=xtls-rprx-vision#MyServer"
	cfg, err := ParseVLESSURI(uri)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.UUID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("UUID: got %q", cfg.UUID)
	}
	if cfg.Address != "example.com" {
		t.Errorf("Address: got %q", cfg.Address)
	}
	if cfg.Port != 443 {
		t.Errorf("Port: got %d", cfg.Port)
	}
	if !cfg.TLS {
		t.Error("TLS should be true")
	}
	if cfg.SNI != "example.com" {
		t.Errorf("SNI: got %q", cfg.SNI)
	}
	if cfg.Transport != config.TransportTCP {
		t.Errorf("Transport: got %q", cfg.Transport)
	}
	if cfg.Flow != "xtls-rprx-vision" {
		t.Errorf("Flow: got %q", cfg.Flow)
	}
	if cfg.Name != "MyServer" {
		t.Errorf("Name: got %q", cfg.Name)
	}
}

func TestParseVLESSURI_WS(t *testing.T) {
	uri := "vless://550e8400-e29b-41d4-a716-446655440000@cdn.example.com:80?type=ws&path=%2Fvless"
	cfg, err := ParseVLESSURI(uri)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Transport != config.TransportWS {
		t.Errorf("Transport: got %q, want ws", cfg.Transport)
	}
	if cfg.Path != "/vless" {
		t.Errorf("Path: got %q, want /vless", cfg.Path)
	}
	if cfg.TLS {
		t.Error("TLS should be false")
	}
}

func TestParseVLESSURI_WSDefaultPath(t *testing.T) {
	uri := "vless://550e8400-e29b-41d4-a716-446655440000@host.com:8080?type=ws"
	cfg, err := ParseVLESSURI(uri)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Path != "/" {
		t.Errorf("Path: got %q, want /", cfg.Path)
	}
}

func TestParseVLESSURI_FallbackName(t *testing.T) {
	uri := "vless://550e8400-e29b-41d4-a716-446655440000@myhost.net:443"
	cfg, err := ParseVLESSURI(uri)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "myhost.net" {
		t.Errorf("Name: got %q, want myhost.net", cfg.Name)
	}
}

func TestParseVLESSURI_InvalidScheme(t *testing.T) {
	_, err := ParseVLESSURI("vmess://uuid@host:443")
	if err == nil {
		t.Error("expected error for unsupported scheme")
	}
}

func TestParseVLESSURI_EmptyUUID(t *testing.T) {
	_, err := ParseVLESSURI("vless://@host:443")
	if err == nil {
		t.Error("expected error for empty UUID")
	}
}

func TestParseVLESSURI_InvalidPort(t *testing.T) {
	_, err := ParseVLESSURI("vless://uuid@host:99999")
	if err == nil {
		t.Error("expected error for invalid port")
	}
}
