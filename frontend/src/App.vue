<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NConfigProvider, NMessageProvider, darkTheme } from 'naive-ui'
import { useAppStore } from './stores/app'
import StatusBar from './components/StatusBar.vue'
import ServerList from './components/ServerList.vue'
import LogsPanel from './components/LogsPanel.vue'

const store = useAppStore()
const activeTab = ref<'servers' | 'logs'>('servers')

onMounted(() => store.init())
</script>

<template>
  <NConfigProvider :theme="darkTheme" class="app-root">
    <NMessageProvider>
      <div class="layout">
        <StatusBar />

        <div class="tab-bar">
          <button
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'servers' }"
            @click="activeTab = 'servers'"
          >
            Серверы
          </button>
          <button
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'logs' }"
            @click="activeTab = 'logs'"
          >
            Логи
            <span v-if="store.logs.some(l => l.level === 'error')" class="tab-badge" />
          </button>
        </div>

        <div class="tab-content">
          <ServerList v-show="activeTab === 'servers'" />
          <LogsPanel v-show="activeTab === 'logs'" />
        </div>
      </div>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style>
.app-root {
  height: 100%;
  background: #1b2636;
}
.layout {
  display: flex;
  flex-direction: column;
  height: 100%;
  color: rgba(255, 255, 255, 0.82);
}
.tab-bar {
  display: flex;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
  padding: 0 16px;
}
.tab-btn {
  position: relative;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.4);
  font-size: 13px;
  padding: 8px 12px;
  cursor: pointer;
  transition: color 0.15s;
}
.tab-btn:hover {
  color: rgba(255, 255, 255, 0.7);
}
.tab-btn--active {
  color: rgba(255, 255, 255, 0.9);
}
.tab-btn--active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: #63e2b7;
  border-radius: 2px 2px 0 0;
}
.tab-badge {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #d03050;
  margin-left: 4px;
  vertical-align: middle;
}
.tab-content {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
.tab-content > * {
  flex: 1;
  min-height: 0;
}
</style>
