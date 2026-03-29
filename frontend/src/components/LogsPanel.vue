<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { NButton, NScrollbar } from 'naive-ui'
import { useAppStore } from '../stores/app'

const store = useAppStore()
const scrollbarRef = ref<InstanceType<typeof NScrollbar> | null>(null)

// авто-скролл вниз при появлении новых записей
watch(
  () => store.logs.length,
  () => nextTick(() => scrollbarRef.value?.scrollTo({ top: 999999, behavior: 'smooth' }))
)

const levelClass: Record<string, string> = {
  error: 'log--error',
  warn: 'log--warn',
  debug: 'log--debug',
}

async function copyLogs() {
  const text = store.logs.map(e => `[${e.time}] [${e.level}] ${e.message}`).join('\n')
  await navigator.clipboard.writeText(text)
}
</script>

<template>
  <div class="logs-panel">
    <div class="logs-toolbar">
      <span class="logs-title">Логи</span>
      <div class="logs-actions">
        <NButton size="tiny" quaternary @click="copyLogs" :disabled="!store.logs.length">
          Копировать
        </NButton>
        <NButton size="tiny" quaternary @click="store.clearLogs" :disabled="!store.logs.length">
          Очистить
        </NButton>
      </div>
    </div>

    <NScrollbar ref="scrollbarRef" class="logs-scroll">
      <div class="logs-list">
        <div v-if="!store.logs.length" class="logs-empty">
          Нет записей
        </div>
        <div
          v-for="(entry, i) in store.logs"
          :key="i"
          class="log-entry"
          :class="levelClass[entry.level]"
        >
          <span class="log-time">{{ entry.time }}</span>
          <span class="log-level">{{ entry.level }}</span>
          <span class="log-msg">{{ entry.message }}</span>
        </div>
      </div>
    </NScrollbar>
  </div>
</template>

<style scoped>
.logs-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.logs-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.logs-title {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.45);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.logs-actions {
  display: flex;
  gap: 4px;
}

.logs-scroll {
  flex: 1;
  min-height: 0;
}

.logs-list {
  padding: 8px 0;
  font-family: 'Menlo', 'Consolas', 'Monaco', monospace;
  font-size: 12px;
}

.logs-empty {
  padding: 24px 16px;
  color: rgba(255, 255, 255, 0.2);
  text-align: center;
  font-family: inherit;
  font-size: 13px;
}

.log-entry {
  display: flex;
  gap: 8px;
  padding: 2px 16px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.75);
}

.log-entry:hover {
  background: rgba(255, 255, 255, 0.03);
}

.log-time {
  color: rgba(255, 255, 255, 0.25);
  flex-shrink: 0;
}

.log-level {
  flex-shrink: 0;
  width: 40px;
  color: rgba(255, 255, 255, 0.3);
}

.log-msg {
  word-break: break-all;
}

/* цвета уровней */
.log--error { color: #e88080; }
.log--error .log-time,
.log--error .log-level { color: #c05050; }

.log--warn { color: #d4b06a; }
.log--warn .log-time,
.log--warn .log-level { color: #a07840; }

.log--debug .log-msg { color: rgba(255, 255, 255, 0.35); }
</style>
