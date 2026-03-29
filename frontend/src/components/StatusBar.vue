<script setup lang="ts">
import { NTag } from 'naive-ui'
import { useAppStore } from '../stores/app'

const store = useAppStore()

const modeLabel: Record<string, string> = {
  proxy: 'Прокси',
  vpn: 'VPN',
  both: 'Прокси + VPN',
}

function formatBytes(n: number): string {
  if (n < 1024) return `${n} Б`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} КБ`
  if (n < 1024 * 1024 * 1024) return `${(n / 1024 / 1024).toFixed(1)} МБ`
  return `${(n / 1024 / 1024 / 1024).toFixed(2)} ГБ`
}
</script>

<template>
  <div class="status-bar">
    <div class="status-left">
      <span class="dot" :class="store.status.connected ? 'dot--on' : 'dot--off'" />
      <span class="status-text">
        {{ store.status.connected ? store.status.server_name : 'Не подключено' }}
      </span>
    </div>
    <div class="status-right">
      <NTag size="small" :bordered="false">{{ modeLabel[store.settings.mode] }}</NTag>
      <template v-if="store.status.connected">
        <span class="traffic">↑ {{ formatBytes(store.stats.upload) }}</span>
        <span class="traffic">↓ {{ formatBytes(store.stats.download) }}</span>
      </template>
      <span class="version">v{{ store.version }}</span>
    </div>
  </div>
</template>

<style scoped>
.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}
.status-left {
  display: flex;
  align-items: center;
  gap: 8px;
}
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot--on { background: #18a058; }
.dot--off { background: #d03050; }
.status-text {
  font-size: 14px;
  font-weight: 500;
}
.status-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.traffic {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
}
.version {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.2);
}
</style>
