<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NModal, NForm, NFormItem, NInput, NInputNumber, NSwitch, NSelect, NSpace, NButton,
  useMessage,
  type FormRules,
  type FormInst,
} from 'naive-ui'
import type { ServerConfig } from '../types'
import { useAppStore } from '../stores/app'

const props = defineProps<{
  visible: boolean
  initial: ServerConfig | null
}>()

const emit = defineEmits<{
  'update:visible': [val: boolean]
  saved: []
}>()

const store = useAppStore()
const message = useMessage()

function defaultForm(): ServerConfig {
  return {
    id: '',
    name: '',
    address: '',
    port: 443,
    uuid: '',
    transport: 'tcp',
    tls: true,
    sni: '',
    path: '/',
    flow: '',
  }
}

const form = ref<ServerConfig>(defaultForm())
const formRef = ref<FormInst | null>(null)
const saving = ref(false)

watch(
  () => props.visible,
  (v) => {
    if (v) {
      if (props.initial) {
        form.value = {
          ...props.initial,
          sni: props.initial.sni ?? '',
          path: props.initial.path ?? '/',
          flow: props.initial.flow ?? '',
        }
      } else {
        form.value = defaultForm()
      }
    }
  },
)

const rules: FormRules = {
  name: [{ required: true, message: 'Введите имя', trigger: 'blur' }],
  address: [{ required: true, message: 'Введите адрес', trigger: 'blur' }],
  port: [{ required: true, type: 'number', message: 'Введите порт', trigger: 'blur' }],
  uuid: [{ required: true, message: 'Введите UUID', trigger: 'blur' }],
}

const transportOptions = [
  { label: 'TCP', value: 'tcp' },
  { label: 'WebSocket', value: 'ws' },
]

async function handleSave() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    await store.saveServer(form.value)
    emit('saved')
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
    :title="initial ? 'Изменить сервер' : 'Добавить сервер'"
    style="width: 480px"
    :auto-focus="false"
    @update:show="emit('update:visible', $event)"
  >
    <NForm ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="110px">
      <NFormItem label="Имя" path="name">
        <NInput v-model:value="form.name" placeholder="Мой сервер" />
      </NFormItem>
      <NFormItem label="Адрес" path="address">
        <NInput v-model:value="form.address" placeholder="example.com" />
      </NFormItem>
      <NFormItem label="Порт" path="port">
        <NInputNumber v-model:value="form.port" :min="1" :max="65535" style="width: 100%" />
      </NFormItem>
      <NFormItem label="UUID" path="uuid">
        <NInput v-model:value="form.uuid" placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" />
      </NFormItem>
      <NFormItem label="Транспорт">
        <NSelect v-model:value="form.transport" :options="transportOptions" />
      </NFormItem>
      <NFormItem label="TLS">
        <NSwitch v-model:value="form.tls" />
      </NFormItem>
      <NFormItem v-if="form.tls" label="SNI">
        <NInput v-model:value="form.sni" placeholder="example.com" />
      </NFormItem>
      <NFormItem v-if="form.transport === 'ws'" label="WS Path">
        <NInput v-model:value="form.path" placeholder="/" />
      </NFormItem>
      <NFormItem label="Flow">
        <NInput v-model:value="form.flow" placeholder="xtls-rprx-vision (необязательно)" />
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
