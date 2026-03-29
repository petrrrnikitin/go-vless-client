export type Mode = 'proxy' | 'vpn' | 'both'
export type Transport = 'tcp' | 'ws'

export interface ServerConfig {
  id: string
  name: string
  address: string
  port: number
  uuid: string
  transport: Transport
  tls: boolean
  sni: string
  path: string
  flow: string
}

export interface AppSettings {
  mode: Mode
  socks5_port: number
  http_port: number
  api_port: number
  last_server_id?: string
}

export interface ConnectionStatus {
  connected: boolean
  server_id?: string
  server_name?: string
  mode: Mode
}

export interface Stats {
  upload: number
  download: number
}