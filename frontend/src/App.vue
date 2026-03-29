<script setup lang="ts">
import { onMounted } from 'vue'
import { NConfigProvider, NMessageProvider, darkTheme } from 'naive-ui'
import { useAppStore } from './stores/app'
import StatusBar from './components/StatusBar.vue'
import ServerList from './components/ServerList.vue'

const store = useAppStore()
onMounted(() => store.init())
</script>

<template>
  <NConfigProvider :theme="darkTheme" class="app-root">
    <NMessageProvider>
      <div class="layout">
        <StatusBar />
        <ServerList />
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
</style>