<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { NSelect, useMessage } from 'naive-ui'
import { useAppStore } from '../stores/app'

const store = useAppStore()
const message = useMessage()

const connecting = ref(false)

const selectedServerId = ref(
  store.status.server_id || store.servers[0]?.id || ''
)

// keep selection in sync when servers load or status changes
watch(() => store.status.server_id, (id) => { if (id) selectedServerId.value = id })
watch(() => store.servers, (list) => {
  if (!selectedServerId.value && list.length) selectedServerId.value = list[0].id
}, { immediate: true })

const serverOptions = computed(() =>
  store.servers.map(s => ({ label: s.name, value: s.id }))
)

const isConnected = computed(() => store.status.connected)
const isCurrentServer = computed(() =>
  store.status.server_id === selectedServerId.value
)

async function toggleConnect() {
  if (isConnected.value) {
    try {
      await store.disconnect()
    } catch (e: unknown) {
      message.error(e instanceof Error ? e.message : 'Ошибка отключения')
    }
    return
  }

  if (!selectedServerId.value) {
    message.warning('Выберите сервер')
    return
  }

  connecting.value = true
  try {
    await store.connect(selectedServerId.value)
  } catch (e: unknown) {
    message.error(e instanceof Error ? e.message : 'Ошибка подключения')
  } finally {
    connecting.value = false
  }
}

async function handleCheckProxy() {
  try {
    const ip = await store.checkProxy()
    message.success(`Внешний IP: ${ip}`)
  } catch (e: unknown) {
    message.error(e instanceof Error ? e.message : 'Ошибка проверки')
  }
}
</script>

<template>
  <div class="home">
    <!-- Server selector -->
    <div class="server-select-wrap">
      <NSelect
        v-model:value="selectedServerId"
        :options="serverOptions"
        placeholder="Выберите сервер"
        :disabled="isConnected"
        class="server-select"
      />
    </div>

    <!-- Connect button -->
    <div class="connect-area">
      <button
        class="connect-btn"
        :class="{
          'connect-btn--connected': isConnected,
          'connect-btn--loading': connecting,
        }"
        :disabled="connecting"
        @click="toggleConnect"
      >
        <!-- Power icon -->
        <svg class="power-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round">
          <path d="M18.36 6.64a9 9 0 1 1-12.73 0" />
          <line x1="12" y1="2" x2="12" y2="12" />
        </svg>
      </button>

      <div class="status-text">
        <template v-if="connecting">Подключение…</template>
        <template v-else-if="isConnected">
          Подключено к <strong>{{ store.status.server_name }}</strong>
        </template>
        <template v-else>Нажмите для подключения</template>
      </div>

      <button
        v-if="isConnected"
        class="check-ip-btn"
        @click="handleCheckProxy"
      >
        Проверить IP
      </button>
    </div>

    <!-- Mode badge -->
    <div class="mode-badge">
      <span>Режим:</span>
      <span class="mode-val">{{ { proxy: 'Прокси', vpn: 'VPN', both: 'Прокси + VPN' }[store.settings.mode] }}</span>
    </div>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  padding: 24px 32px;
  gap: 0;
}

.server-select-wrap {
  width: 100%;
  max-width: 340px;
}

.server-select {
  width: 100%;
}

.connect-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.connect-btn {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 2.5px solid rgba(255,255,255,0.12);
  background: rgba(255,255,255,0.04);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.25s ease;
  color: rgba(255,255,255,0.35);
  position: relative;
}

.connect-btn::before {
  content: '';
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  border: 1px solid transparent;
  transition: all 0.25s ease;
}

.connect-btn:not(:disabled):hover {
  border-color: rgba(99, 226, 183, 0.4);
  background: rgba(99, 226, 183, 0.06);
  color: rgba(99, 226, 183, 0.7);
}

.connect-btn:not(:disabled):hover::before {
  border-color: rgba(99, 226, 183, 0.15);
}

.connect-btn--connected {
  border-color: #18a058;
  background: rgba(24, 160, 88, 0.12);
  color: #18a058;
  box-shadow: 0 0 24px rgba(24, 160, 88, 0.2);
}

.connect-btn--connected::before {
  border-color: rgba(24, 160, 88, 0.2);
}

.connect-btn--connected:hover {
  border-color: #d03050 !important;
  background: rgba(208, 48, 80, 0.1) !important;
  color: #d03050 !important;
  box-shadow: none;
}

.connect-btn--loading {
  opacity: 0.6;
  cursor: not-allowed;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 0.3; }
}

.power-icon {
  width: 44px;
  height: 44px;
}

.status-text {
  font-size: 14px;
  color: rgba(255,255,255,0.5);
  text-align: center;
}

.status-text strong {
  color: rgba(255,255,255,0.85);
  font-weight: 500;
}

.check-ip-btn {
  background: none;
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 6px;
  color: rgba(255,255,255,0.4);
  font-size: 12px;
  padding: 4px 12px;
  cursor: pointer;
  transition: all 0.15s;
}

.check-ip-btn:hover {
  border-color: rgba(255,255,255,0.25);
  color: rgba(255,255,255,0.7);
}

.mode-badge {
  font-size: 12px;
  color: rgba(255,255,255,0.2);
  display: flex;
  gap: 6px;
}

.mode-val {
  color: rgba(255,255,255,0.4);
}
</style>
