<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NModal, NForm, NFormItem, NInputNumber, NSelect, NSpace, NButton,
  useMessage,
} from 'naive-ui'
import type { AppSettings } from '../types'
import { useAppStore } from '../stores/app'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  'update:visible': [val: boolean]
}>()

const store = useAppStore()
const message = useMessage()

const form = ref<AppSettings>({ ...store.settings })
const saving = ref(false)

watch(
  () => props.visible,
  (v) => {
    if (v) form.value = { ...store.settings }
  },
)

const modeOptions = [
  { label: 'Только прокси (SOCKS5/HTTP)', value: 'proxy' },
  { label: 'Только VPN (TUN)', value: 'vpn' },
  { label: 'Прокси + VPN', value: 'both' },
]

async function handleSave() {
  saving.value = true
  try {
    await store.saveSettings(form.value)
    emit('update:visible', false)
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : 'Ошибка сохранения'
    message.error(msg)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <NModal
    :show="visible"
    preset="card"
    title="Настройки"
    style="width: 420px"
    :auto-focus="false"
    @update:show="emit('update:visible', $event)"
  >
    <NForm :model="form" label-placement="left" label-width="130px">
      <NFormItem label="Режим">
        <NSelect v-model:value="form.mode" :options="modeOptions" />
      </NFormItem>
      <NFormItem label="SOCKS5 порт">
        <NInputNumber v-model:value="form.socks5_port" :min="1" :max="65535" style="width: 100%" />
      </NFormItem>
      <NFormItem label="HTTP порт">
        <NInputNumber v-model:value="form.http_port" :min="1" :max="65535" style="width: 100%" />
      </NFormItem>
      <NFormItem label="API порт">
        <NInputNumber v-model:value="form.api_port" :min="1" :max="65535" style="width: 100%" />
      </NFormItem>
    </NForm>
    <template #footer>
      <NSpace justify="end">
        <NButton @click="emit('update:visible', false)">Отмена</NButton>
        <NButton type="primary" :loading="saving" @click="handleSave">Сохранить</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
