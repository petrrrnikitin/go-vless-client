<script setup lang="ts">
import { ref, watch } from 'vue'
import { NModal, NButton, NInput, NSpace, NText, useMessage } from 'naive-ui'
import { ParseURI } from '../../wailsjs/go/main/App'
import type { ServerConfig } from '../types'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  'update:visible': [val: boolean]
  parsed: [cfg: ServerConfig]
}>()

const message = useMessage()
const uri = ref('')
const loading = ref(false)
const fromClipboard = ref(false)

watch(
  () => props.visible,
  async (v) => {
    if (!v) return
    uri.value = ''
    fromClipboard.value = false
    try {
      const text = await navigator.clipboard.readText()
      if (text.trimStart().startsWith('vless://')) {
        uri.value = text.trim()
        fromClipboard.value = true
      }
    } catch {
      // clipboard access denied — ignore
    }
  },
)

async function handleImport() {
  if (!uri.value.trim()) return
  loading.value = true
  try {
    const cfg = await ParseURI(uri.value.trim())
    emit('parsed', cfg)
    emit('update:visible', false)
    uri.value = ''
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : 'Неверный URI'
    message.error(msg)
  } finally {
    loading.value = false
  }
}

function handleClose() {
  uri.value = ''
  emit('update:visible', false)
}
</script>

<template>
  <NModal
    :show="visible"
    preset="card"
    title="Импорт из VLESS URI"
    style="width: 520px"
    :auto-focus="false"
    @update:show="emit('update:visible', $event)"
  >
    <NText v-if="fromClipboard" depth="3" style="display: block; margin-bottom: 8px; font-size: 12px">
      Обнаружен URI из буфера обмена
    </NText>
    <NInput
      v-model:value="uri"
      type="textarea"
      placeholder="vless://uuid@host:port?security=tls&type=ws#Название"
      :rows="3"
      @keydown.enter.ctrl="handleImport"
    />
    <template #footer>
      <NSpace justify="end">
        <NButton @click="handleClose">Отмена</NButton>
        <NButton type="primary" :loading="loading" @click="handleImport">Импортировать</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
