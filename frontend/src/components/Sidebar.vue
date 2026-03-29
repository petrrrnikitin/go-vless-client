<script setup lang="ts">
import { useAppStore } from '../stores/app'
import SettingsModal from './SettingsModal.vue'
import { ref } from 'vue'

const props = defineProps<{ active: string }>()
const emit = defineEmits<{ navigate: [page: string] }>()

const store = useAppStore()
const showSettings = ref(false)

function formatBytes(n: number): string {
  if (n < 1024) return `${n} Б`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} КБ`
  if (n < 1024 * 1024 * 1024) return `${(n / 1024 / 1024).toFixed(1)} МБ`
  return `${(n / 1024 / 1024 / 1024).toFixed(2)} ГБ`
}

const navItems = [
  { id: 'home',    label: 'Главная' },
  { id: 'servers', label: 'Серверы' },
  { id: 'logs',    label: 'Логи' },
]
</script>

<template>
  <div class="sidebar">
    <!-- App header -->
    <div class="sidebar-header">
      <span class="app-title">VPN Client</span>
      <span class="version-badge">v{{ store.version }}</span>
    </div>

    <!-- Navigation -->
    <nav class="nav">
      <button
        v-for="item in navItems"
        :key="item.id"
        class="nav-item"
        :class="{ 'nav-item--active': active === item.id }"
        @click="emit('navigate', item.id)"
      >
        <span class="nav-label">{{ item.label }}</span>
        <span
          v-if="item.id === 'logs' && store.logs.some(l => l.level === 'error')"
          class="nav-dot"
        />
      </button>
    </nav>

    <div class="sidebar-spacer" />

    <!-- Traffic stats -->
    <div class="traffic" v-if="store.status.connected">
      <div class="traffic-title">Трафик</div>
      <div class="traffic-row">
        <span class="traffic-arrow up">↑</span>
        <span class="traffic-val">{{ formatBytes(store.stats.upload) }}</span>
      </div>
      <div class="traffic-row">
        <span class="traffic-arrow down">↓</span>
        <span class="traffic-val">{{ formatBytes(store.stats.download) }}</span>
      </div>
    </div>

    <!-- Settings -->
    <button class="nav-item nav-item--settings" @click="showSettings = true">
      <span class="nav-label">Настройки</span>
    </button>
  </div>

  <SettingsModal v-model:visible="showSettings" />
</template>

<style scoped>
.sidebar {
  width: 192px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: #141d2b;
  border-right: 1px solid rgba(255,255,255,0.07);
  padding: 0 0 8px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px 16px 16px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  margin-bottom: 8px;
}

.app-title {
  font-size: 15px;
  font-weight: 600;
  color: rgba(255,255,255,0.9);
  letter-spacing: 0.01em;
}

.version-badge {
  font-size: 11px;
  color: rgba(255,255,255,0.3);
  background: rgba(255,255,255,0.06);
  padding: 1px 6px;
  border-radius: 10px;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 0 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 9px 10px;
  border-radius: 6px;
  background: none;
  border: none;
  color: rgba(255,255,255,0.45);
  font-size: 13.5px;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  text-align: left;
}

.nav-item:hover {
  background: rgba(255,255,255,0.05);
  color: rgba(255,255,255,0.75);
}

.nav-item--active {
  background: rgba(99, 226, 183, 0.1);
  color: #63e2b7;
}

.nav-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #d03050;
  flex-shrink: 0;
}

.sidebar-spacer {
  flex: 1;
}

.traffic {
  padding: 12px 16px;
  margin: 0 8px;
  border-radius: 6px;
  background: rgba(255,255,255,0.03);
  margin-bottom: 8px;
}

.traffic-title {
  font-size: 11px;
  color: rgba(255,255,255,0.25);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 6px;
}

.traffic-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 2px 0;
}

.traffic-arrow {
  font-size: 12px;
  width: 12px;
}
.traffic-arrow.up   { color: #63e2b7; }
.traffic-arrow.down { color: #70b8ff; }

.traffic-val {
  font-size: 12px;
  color: rgba(255,255,255,0.55);
}

.nav-item--settings {
  margin: 0 8px;
  width: calc(100% - 16px);
}
</style>
