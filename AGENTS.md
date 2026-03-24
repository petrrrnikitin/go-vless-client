# AGENTS.md

## Описание проекта

**go-vless-client** — десктопное приложение для подключения к VLESS-прокси серверам с GUI на основе Wails.

### Цель

Предоставить полноценный VPN-клиент с двумя режимами работы:

- **Режим прокси** — локальный SOCKS5 и HTTP прокси. Работает без root. Используется браузером и приложениями, поддерживающими прокси.
- **Режим VPN (TUN)** — создаёт виртуальный сетевой интерфейс, весь трафик системы проходит через туннель. Требует root / `CAP_NET_ADMIN`.

### Сопутствующий проект

Браузерное расширение (отдельный репозиторий), которое взаимодействует с приложением через локальный REST API для переключения прокси в браузере.

---

## Технологический стек

| Компонент | Технология |
|---|---|
| Desktop GUI | [Wails v2](https://wails.io/) |
| Backend | Go 1.22+ |
| VLESS / транспорт / TUN | [sing-box](https://github.com/sagernet/sing-box) |
| Локальный SOCKS5 прокси | `localhost:1080` (настраивается) |
| Локальный HTTP прокси | `localhost:8080` (настраивается) |
| REST API для расширения | `localhost:9090` (настраивается) |
| Frontend | Vue 3 + TypeScript (Composition API, `<script setup>`) |
| UI компоненты | Naive UI |
| Конфигурация | JSON в `os.UserConfigDir()/go-vless-client/` |

---

## Режимы работы

### Режим прокси

- Запускает локальный SOCKS5 (`localhost:1080`) и HTTP (`localhost:8080`) прокси
- Не требует привилегий
- Трафик проксируется только у приложений, настроенных на этот прокси
- Используется браузерным расширением

### Режим VPN (TUN)

- Создаёт TUN-интерфейс через sing-box
- Весь сетевой трафик системы перенаправляется в туннель
- Встроенный DNS для предотвращения DNS-утечек
- Требует root / `CAP_NET_ADMIN` на Linux, Administrator на Windows
- При включении запрашивает повышение прав

Оба режима могут быть активны одновременно.

---

## Архитектура

```
Browser Extension
      │ HTTP (localhost:9090)
      ▼
┌──────────────────────────────────────────┐
│               Wails App                  │
│  ┌──────────┐  ┌────────────────────┐   │
│  │ Vue 3 UI │  │     REST API       │   │
│  └────┬─────┘  └─────────┬──────────┘   │
│       │   Wails bridge   │              │
│  ┌────▼──────────────────▼──────────┐   │
│  │          App struct (Go)         │   │
│  │  Connect / Disconnect            │   │
│  │  SetMode (proxy | vpn | both)    │   │
│  │  GetServers / SaveServer         │   │
│  │  GetStatus / GetStats            │   │
│  └──────────────┬───────────────────┘   │
│                 │                        │
│  ┌──────────────▼───────────────────┐   │
│  │         sing-box core            │   │
│  │  VLESS + TCP/WS/TLS              │   │
│  │  ┌─────────────┐ ┌────────────┐  │   │
│  │  │  TUN mode   │ │ Proxy mode │  │   │
│  │  │ (TUN iface) │ │ SOCKS5/HTTP│  │   │
│  │  └─────────────┘ └────────────┘  │   │
│  └──────────────┬───────────────────┘   │
└─────────────────┼────────────────────── ┘
                  │ туннель
           VLESS Server
```

---

## Структура проекта

```
go-vless-client/
├── AGENTS.md
├── main.go                        # Wails entry point
├── app.go                         # App struct — методы для frontend
├── go.mod
├── go.sum
├── wails.json
│
├── internal/
│   ├── config/
│   │   ├── model.go               # ServerConfig, AppSettings, Mode
│   │   └── storage.go             # Загрузка/сохранение JSON
│   ├── core/
│   │   ├── client.go              # Обёртка над sing-box (proxy + TUN)
│   │   └── stats.go               # Счётчики трафика
│   └── api/
│       └── server.go              # REST API для браузерного расширения
│
├── build/                         # Wails build assets (иконки, манифест)
│
└── frontend/
    ├── src/
    │   ├── App.vue
    │   ├── main.ts
    │   ├── components/
    │   │   ├── ServerList.vue      # Список серверов
    │   │   ├── ServerForm.vue      # Форма добавления/редактирования
    │   │   ├── StatusBar.vue       # Статус + трафик
    │   │   ├── ModeSwitch.vue      # Переключатель режима: прокси / VPN / оба
    │   │   └── SettingsModal.vue   # Настройки портов и т.д.
    │   ├── stores/
    │   │   ├── connection.ts       # Pinia: статус подключения и режим
    │   │   └── servers.ts          # Pinia: список серверов
    │   ├── types/
    │   │   └── index.ts            # ServerConfig, ConnectionStatus, Stats, Mode
    │   └── wailsjs/                # Авто-генерируется Wails (не редактировать)
    ├── index.html
    ├── vite.config.ts
    ├── tsconfig.json
    └── package.json
```

---

## MVP Scope

Что входит в MVP:

- [ ] Добавление/удаление серверов (UUID, адрес, порт, транспорт)
- [ ] Импорт сервера по `vless://` URI
- [ ] Режим прокси: SOCKS5 + HTTP без root
- [ ] Режим VPN (TUN): весь трафик системы через туннель
- [ ] Переключение режима в UI
- [ ] Подключение и отключение
- [ ] Статус подключения в UI
- [ ] Счётчик трафика (↑ / ↓)
- [ ] REST API: `GET /status` → `{ connected, mode, server, proxies }`
- [ ] Сохранение конфигов между перезапусками

Что **не входит** в MVP:

- Браузерное расширение
- System tray
- Автоподключение при старте
- Split tunneling (роутинг по правилам)
- Kill switch
- Обновление приложения

---

## Транспорты (MVP)

| Транспорт | Поддержка |
|---|---|
| TCP | ✅ MVP |
| WebSocket | ✅ MVP |
| TLS | ✅ MVP |
| gRPC | ⬜ После MVP |
| HTTP/2 | ⬜ После MVP |

---

## Wails App Struct

Публичные методы `App`, доступные из Vue через `wailsjs`:

```go
// Серверы
GetServers() []config.ServerConfig
SaveServer(cfg config.ServerConfig) error
DeleteServer(id string) error

// Подключение
Connect(serverID string) error
Disconnect() error
GetStatus() core.ConnectionStatus

// Режим работы
SetMode(mode config.Mode) error  // "proxy" | "vpn" | "both"
GetMode() config.Mode

// Настройки
GetSettings() config.AppSettings
SaveSettings(s config.AppSettings) error
```

Events (Go → Vue через `runtime.EventsEmit`):

| Событие | Payload | Описание |
|---|---|---|
| `status:changed` | `ConnectionStatus` | Смена статуса подключения |
| `stats:update` | `Stats` | Обновление трафика (каждые 2с) |
| `mode:changed` | `Mode` | Смена режима работы |
| `error` | `string` | Ошибка (например, нет прав для TUN) |

---

## REST API для браузерного расширения

Сервер слушает на `localhost:9090` (только loopback).

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/status` | Статус подключения, режим и порты прокси |
| `POST` | `/connect` | Подключиться к последнему серверу |
| `POST` | `/disconnect` | Отключиться |

Пример ответа `GET /status`:
```json
{
  "connected": true,
  "mode": "both",
  "server": "my-server",
  "http_proxy": "localhost:8080",
  "socks5_proxy": "localhost:1080"
}
```

---

## Соглашения по коду

### Go

- Форматирование: `gofmt` / `goimports`
- Линтер: `golangci-lint`
- Именование: `camelCase` для неэкспортируемых, `PascalCase` для экспортируемых
- Ошибки оборачиваются: `fmt.Errorf("connect: %w", err)`
- Интерфейсы определяются там, где **используются**, не где реализованы
- Контекст (`context.Context`) — первый аргумент функций, которые могут блокироваться
- Без `init()`, без глобальных переменных (кроме ошибок-sentinel)

### Vue 3 / TypeScript

- Composition API, `<script setup lang="ts">` везде
- Pinia для глобального состояния
- Компоненты именуются `PascalCase.vue`
- Props и emits явно типизированы через `defineProps<>` / `defineEmits<>`
- Без `any`, без `// @ts-ignore`

### Коммиты

Conventional Commits:

```
feat: добавить импорт vless:// URI
fix: исправить утечку горутины при дисконнекте
chore: обновить зависимости
refactor: вынести логику прокси в отдельный пакет
```

---

## Локальная разработка

```bash
# Установить Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Запустить в dev-режиме (hot reload фронта)
wails dev

# Собрать бинарь
wails build

# TUN-режим требует прав при запуске
sudo ./build/bin/go-vless-client
```

---

## Зависимости (ключевые)

```
github.com/wailsapp/wails/v2
github.com/sagernet/sing-box
github.com/google/uuid
```