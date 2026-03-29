<script setup lang="ts">
import { ref } from 'vue'
import { NEmpty, NTag, NButton, NSpace, NPopconfirm } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { useAppStore } from '../stores/app'
import type { ServerConfig } from '../types'
import ServerForm from './ServerForm.vue'
import SettingsModal from './SettingsModal.vue'
import URIImport from './URIImport.vue'

const store = useAppStore()
const message = useMessage()

const showForm = ref(false)
const editServer = ref<ServerConfig | null>(null)
const showSettings = ref(false)
const showURIImport = ref(false)

function openURIImport() {
  showURIImport.value = true
}

function handleParsed(cfg: ServerConfig) {
  editServer.value = cfg
  showForm.value = true
}
const pingResults = ref<Record<string, string>>({})
const connecting = ref<string | null>(null)

function openAdd() {
  editServer.value = null
  showForm.value = true
}

function openEdit(srv: ServerConfig) {
  editServer.value = srv
  showForm.value = true
}

async function handleConnect(id: string) {
  connecting.value = id
  try {
    await store.connect(id)
  } catch (e: unknown) {
    message.error(errorMessage(e))
  } finally {
    connecting.value = null
  }
}

async function handleDisconnect() {
  try {
    await store.disconnect()
  } catch (e: unknown) {
    message.error(errorMessage(e))
  }
}

async function handlePing(srv: ServerConfig) {
  pingResults.value[srv.id] = '...'
  try {
    const ms = await store.ping(srv.id)
    pingResults.value[srv.id] = `${ms} мс`
  } catch {
    pingResults.value[srv.id] = 'ошибка'
  }
}

async function handleDelete(id: string) {
  try {
    await store.deleteServer(id)
  } catch (e: unknown) {
    message.error(errorMessage(e))
  }
}

async function handleCheckProxy() {
  try {
    const ip = await store.checkProxy()
    message.success(`Внешний IP: ${ip}`)
  } catch (e: unknown) {
    message.error(errorMessage(e))
  }
}

function errorMessage(e: unknown): string {
  if (e instanceof Error) return e.message
  if (typeof e === 'string') return e
  return 'Неизвестная ошибка'
}
</script>

<template>
  <div class="server-list">
    <div class="server-items">
      <NEmpty v-if="!store.servers.length" description="Нет серверов" class="empty" />

      <div
        v-for="srv in store.servers"
        :key="srv.id"
        class="server-card"
        :class="{ 'server-card--active': store.status.server_id === srv.id }"
      >
        <div class="server-info">
          <span class="server-name">{{ srv.name }}</span>
          <span class="server-addr">{{ srv.address }}:{{ srv.port }}</span>
          <div class="server-tags">
            <NTag size="small" :bordered="false">{{ srv.transport.toUpperCase() }}</NTag>
            <NTag v-if="srv.tls" size="small" type="info" :bordered="false">TLS</NTag>
          </div>
          <span v-if="pingResults[srv.id]" class="ping-result">{{ pingResults[srv.id] }}</span>
        </div>
        <div class="server-actions">
          <NButton size="small" @click="handlePing(srv)">Пинг</NButton>
          <NButton size="small" @click="openEdit(srv)">Изменить</NButton>
          <NPopconfirm @positive-click="handleDelete(srv.id)">
            <template #trigger>
              <NButton size="small" type="error" ghost>Удалить</NButton>
            </template>
            Удалить сервер?
          </NPopconfirm>
          <NButton
            v-if="store.status.server_id !== srv.id"
            size="small"
            type="primary"
            :loading="connecting === srv.id"
            :disabled="store.status.connected && store.status.server_id !== srv.id"
            @click="handleConnect(srv.id)"
          >Подключить</NButton>
          <NButton
            v-else
            size="small"
            type="warning"
            @click="handleDisconnect"
          >Отключить</NButton>
        </div>
      </div>
    </div>

    <div class="bottom-bar">
      <NSpace>
        <NButton type="primary" @click="openAdd">+ Добавить сервер</NButton>
        <NButton @click="openURIImport">Импортировать URI</NButton>
        <NButton v-if="store.status.connected" @click="handleCheckProxy">Проверить IP</NButton>
        <NButton @click="showSettings = true">Настройки</NButton>
      </NSpace>
    </div>

    <ServerForm
      v-model:visible="showForm"
      :initial="editServer"
      @saved="showForm = false"
    />
    <SettingsModal v-model:visible="showSettings" />
    <URIImport v-model:visible="showURIImport" @parsed="handleParsed" />
  </div>
</template>

<style scoped>
.server-list {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}
.server-items {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.empty {
  margin-top: 60px;
}
.server-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: border-color 0.2s;
  gap: 12px;
}
.server-card--active {
  border-color: #18a058;
}
.server-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}
.server-name {
  font-weight: 500;
  font-size: 14px;
}
.server-addr {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
}
.server-tags {
  display: flex;
  gap: 4px;
  margin-top: 2px;
}
.ping-result {
  font-size: 11px;
  color: #18a058;
  margin-top: 2px;
}
.server-actions {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-shrink: 0;
}
.bottom-bar {
  padding: 12px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}
</style>
