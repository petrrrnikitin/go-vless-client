<script setup lang="ts">
import { ref, h } from 'vue'
import { NEmpty, NTag, NButton, NPopconfirm, NDropdown, useMessage } from 'naive-ui'
import { ExportURI } from '../../wailsjs/go/main/App'
import { useAppStore } from '../stores/app'
import type { ServerConfig } from '../types'
import ServerForm from './ServerForm.vue'
import URIImport from './URIImport.vue'

const store = useAppStore()
const message = useMessage()

const showForm = ref(false)
const editServer = ref<ServerConfig | null>(null)
const showURIImport = ref(false)
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

function handleParsed(cfg: ServerConfig) {
  editServer.value = cfg
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
  pingResults.value[srv.id] = '…'
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

async function handleCopyURI(id: string) {
  try {
    const uri = await ExportURI(id)
    await navigator.clipboard.writeText(uri)
    message.success('URI скопирован')
  } catch (e: unknown) {
    message.error(errorMessage(e))
  }
}

function errorMessage(e: unknown): string {
  if (e instanceof Error) return e.message
  if (typeof e === 'string') return e
  return 'Неизвестная ошибка'
}

// "⋯" dropdown per server
function menuOptions(srv: ServerConfig) {
  return [
    { label: 'Изменить',      key: 'edit' },
    { label: 'Копировать URI', key: 'copy' },
    { label: 'Удалить',       key: 'delete' },
  ]
}

function handleMenuSelect(key: string, srv: ServerConfig) {
  if (key === 'edit')   openEdit(srv)
  if (key === 'copy')   handleCopyURI(srv.id)
  if (key === 'delete') handleDelete(srv.id)
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
        <!-- Info -->
        <div class="server-info">
          <div class="server-name-row">
            <span class="server-name">{{ srv.name }}</span>
            <span v-if="pingResults[srv.id]" class="ping-result">{{ pingResults[srv.id] }}</span>
          </div>
          <span class="server-addr">{{ srv.address }}:{{ srv.port }}</span>
          <div class="server-tags">
            <NTag size="small" :bordered="false">{{ srv.transport.toUpperCase() }}</NTag>
            <NTag v-if="srv.tls" size="small" type="info" :bordered="false">TLS</NTag>
          </div>
        </div>

        <!-- Actions -->
        <div class="server-actions">
          <NButton size="small" @click="handlePing(srv)">Пинг</NButton>

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

          <NDropdown
            :options="menuOptions(srv)"
            @select="(key: string) => handleMenuSelect(key, srv)"
          >
            <NButton size="small" quaternary>⋯</NButton>
          </NDropdown>
        </div>
      </div>
    </div>

    <!-- Bottom bar -->
    <div class="bottom-bar">
      <NButton type="primary" @click="openAdd">+ Добавить</NButton>
      <NButton @click="showURIImport = true">Импортировать URI</NButton>
    </div>

    <ServerForm
      v-model:visible="showForm"
      :initial="editServer"
      @saved="showForm = false"
    />
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
  gap: 6px;
}

.empty {
  margin-top: 60px;
}

.server-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: border-color 0.2s;
  gap: 12px;
}

.server-card--active {
  border-color: #18a058;
  background: rgba(24, 160, 88, 0.05);
}

.server-info {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
  flex: 1;
}

.server-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.server-name {
  font-weight: 500;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ping-result {
  font-size: 11px;
  color: #18a058;
  flex-shrink: 0;
}

.server-addr {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

.server-tags {
  display: flex;
  gap: 4px;
  margin-top: 2px;
}

.server-actions {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-shrink: 0;
}

.bottom-bar {
  display: flex;
  gap: 8px;
  padding: 10px 14px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}
</style>
