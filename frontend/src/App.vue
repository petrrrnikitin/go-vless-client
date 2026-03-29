<script setup lang="ts">
import { ref } from 'vue'
import { NConfigProvider, NMessageProvider, darkTheme } from 'naive-ui'
import { useAppStore } from './stores/app'
import Sidebar from './components/Sidebar.vue'
import HomePage from './components/HomePage.vue'
import ServerList from './components/ServerList.vue'
import LogsPanel from './components/LogsPanel.vue'

const store = useAppStore()
const page = ref<'home' | 'servers' | 'logs'>('home')
</script>

<template>
  <NConfigProvider :theme="darkTheme" class="app-root">
    <NMessageProvider>
      <div class="layout">
        <Sidebar :active="page" @navigate="page = $event as typeof page" />
        <main class="content">
          <HomePage   v-if="page === 'home'" />
          <ServerList v-else-if="page === 'servers'" />
          <LogsPanel  v-else-if="page === 'logs'" />
        </main>
      </div>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }

body, html, #app {
  height: 100%;
  overflow: hidden;
}

.app-root {
  height: 100%;
  background: #1b2636;
}

.layout {
  display: flex;
  height: 100%;
  color: rgba(255, 255, 255, 0.82);
}

.content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
