<script setup lang="ts">
import { ref } from 'vue'
import { NModal, NButton, NInput, NSpace, useMessage } from 'naive-ui'
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
