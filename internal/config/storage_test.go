package config

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStorage(t *testing.T) *Storage {
	t.Helper()
	dir := t.TempDir()
	s := &Storage{
		path: filepath.Join(dir, "config.json"),
		data: configData{Settings: DefaultSettings()},
	}
	return s
}

func TestDefaultSettings(t *testing.T) {
	s := newTestStorage(t)
	settings := s.Settings()

	if settings.Mode != ModeProxy {
		t.Errorf("default mode = %q, want %q", settings.Mode, ModeProxy)
	}
	if settings.Socks5Port != 1080 {
		t.Errorf("default socks5 port = %d, want 1080", settings.Socks5Port)
	}
	if settings.HTTPPort != 8080 {
		t.Errorf("default http port = %d, want 8080", settings.HTTPPort)
	}
}

func TestSaveAndLoadServer(t *testing.T) {
	s := newTestStorage(t)

	srv := ServerConfig{
		ID:        "test-id",
		Name:      "Test Server",
		Address:   "example.com",
		Port:      443,
		UUID:      "some-uuid",
		Transport: TransportWS,
		TLS:       true,
	}

	if err := s.SaveServer(srv); err != nil {
		t.Fatalf("SaveServer: %v", err)
	}

	// перезагружаем с диска
	s2 := &Storage{path: s.path, data: configData{Settings: DefaultSettings()}}
	if err := s2.load(); err != nil {
		t.Fatalf("load: %v", err)
	}

	got, ok := s2.ServerByID("test-id")
	if !ok {
		t.Fatal("server not found after reload")
	}
	if got.Name != srv.Name {
		t.Errorf("name = %q, want %q", got.Name, srv.Name)
	}
	if got.Address != srv.Address {
		t.Errorf("address = %q, want %q", got.Address, srv.Address)
	}
}

func TestUpdateServer(t *testing.T) {
	s := newTestStorage(t)

	srv := ServerConfig{ID: "id-1", Name: "Old Name", Address: "old.example.com", Port: 443, UUID: "uuid"}
	if err := s.SaveServer(srv); err != nil {
		t.Fatalf("SaveServer: %v", err)
	}

	srv.Name = "New Name"
	if err := s.SaveServer(srv); err != nil {
		t.Fatalf("SaveServer update: %v", err)
	}

	if len(s.Servers()) != 1 {
		t.Errorf("expected 1 server, got %d", len(s.Servers()))
	}
	if s.Servers()[0].Name != "New Name" {
		t.Errorf("name = %q, want %q", s.Servers()[0].Name, "New Name")
	}
}

func TestDeleteServer(t *testing.T) {
	s := newTestStorage(t)

	if err := s.SaveServer(ServerConfig{ID: "id-1", Name: "Server 1", UUID: "uuid1"}); err != nil {
		t.Fatalf("SaveServer id-1: %v", err)
	}
	if err := s.SaveServer(ServerConfig{ID: "id-2", Name: "Server 2", UUID: "uuid2"}); err != nil {
		t.Fatalf("SaveServer id-2: %v", err)
	}

	if err := s.DeleteServer("id-1"); err != nil {
		t.Fatalf("DeleteServer: %v", err)
	}

	servers := s.Servers()
	if len(servers) != 1 {
		t.Errorf("expected 1 server after delete, got %d", len(servers))
	}
	if servers[0].ID != "id-2" {
		t.Errorf("remaining server id = %q, want id-2", servers[0].ID)
	}
}

func TestSaveSettings(t *testing.T) {
	s := newTestStorage(t)

	settings := AppSettings{
		Mode:       ModeVPN,
		Socks5Port: 1081,
		HTTPPort:   8081,
		APIPort:    9091,
	}

	if err := s.SaveSettings(settings); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	s2 := &Storage{path: s.path, data: configData{Settings: DefaultSettings()}}
	if err := s2.load(); err != nil {
		t.Fatalf("load: %v", err)
	}

	got := s2.Settings()
	if got.Mode != ModeVPN {
		t.Errorf("mode = %q, want %q", got.Mode, ModeVPN)
	}
	if got.Socks5Port != 1081 {
		t.Errorf("socks5 port = %d, want 1081", got.Socks5Port)
	}
}

func TestAtomicSave(t *testing.T) {
	s := newTestStorage(t)
	if err := s.SaveServer(ServerConfig{ID: "id-1", Name: "Server", UUID: "uuid"}); err != nil {
		t.Fatalf("SaveServer: %v", err)
	}

	// tmp файлов не должно остаться
	entries, _ := os.ReadDir(filepath.Dir(s.path))
	for _, e := range entries {
		if e.Name() != "config.json" {
			t.Errorf("unexpected file after save: %s", e.Name())
		}
	}
}