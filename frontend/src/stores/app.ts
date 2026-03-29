import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ServerConfig, AppSettings, ConnectionStatus, Stats } from '../types'
import * as App from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export const useAppStore = defineStore('app', () => {
  const servers = ref<ServerConfig[]>([])
  const settings = ref<AppSettings>({
    mode: 'proxy',
    socks5_port: 1080,
    http_port: 8080,
    api_port: 9090,
    auto_connect: false,
  })
  const status = ref<ConnectionStatus>({ connected: false, mode: 'proxy' })
  const stats = ref<Stats>({ upload: 0, download: 0 })

  async function init() {
    const [srv, stg, st] = await Promise.all([
      App.GetServers(),
      App.GetSettings(),
      App.GetStatus(),
    ])
    servers.value = srv ?? []
    settings.value = stg
    status.value = st

    EventsOn('status:changed', (data: ConnectionStatus) => {
      status.value = data
      if (!data.connected) stats.value = { upload: 0, download: 0 }
    })
    EventsOn('stats:update', (data: Stats) => {
      stats.value = data
    })
  }

  async function loadServers() {
    servers.value = (await App.GetServers()) ?? []
  }

  async function saveServer(cfg: ServerConfig) {
    await App.SaveServer(cfg)
    await loadServers()
  }

  async function deleteServer(id: string) {
    await App.DeleteServer(id)
    await loadServers()
  }

  async function connect(serverID: string) {
    await App.Connect(serverID)
  }

  async function disconnect() {
    await App.Disconnect()
    stats.value = { upload: 0, download: 0 }
  }

  async function saveSettings(s: AppSettings) {
    await App.SaveSettings(s)
    settings.value = s
  }

  async function ping(serverID: string): Promise<number> {
    return App.Ping(serverID)
  }

  async function checkProxy(): Promise<string> {
    return App.CheckProxy()
  }

  return {
    servers, settings, status, stats,
    init, loadServers, saveServer, deleteServer,
    connect, disconnect, saveSettings, ping, checkProxy,
  }
})