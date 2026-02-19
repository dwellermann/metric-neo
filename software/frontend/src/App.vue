<script setup>
import { ref, onMounted } from 'vue';
import { darkTheme, NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui';
import { NeedsSetup, GetSystemTheme } from '../wailsjs/go/main/App';
import SetupDialog from './components/SetupDialog.vue';
import MainApp from './components/MainApp.vue';

const needsSetup = ref(true);
const isLoading = ref(true);
const naiveTheme = ref(null);

// Apply theme from Go backend
async function applyThemeFromSystem() {
  try {
    const systemTheme = await GetSystemTheme();
    console.log('System theme from Go backend:', systemTheme);

    naiveTheme.value = systemTheme === 'dark' ? darkTheme : null;
  } catch (error) {
    console.error('Failed to get system theme:', error);
    // Fallback zu Browser detection
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    naiveTheme.value = prefersDark ? darkTheme : null;
  }
}

onMounted(async () => {
  // Apply theme on startup
  await applyThemeFromSystem();

  try {
    // Prüfe ob Setup benötigt wird
    needsSetup.value = await NeedsSetup();
    console.log('Setup needed:', needsSetup.value);
  } catch (error) {
    console.error('Failed to check setup status:', error);
  } finally {
    isLoading.value = false;
  }
});

function onSetupComplete(dataDir) {
  console.log('Setup completed! Data directory:', dataDir);
  needsSetup.value = false;
}
</script>

<template>
  <n-config-provider :theme="naiveTheme">
    <n-message-provider>
      <n-dialog-provider>
        <!-- Ladebildschirm -->
        <div v-if="isLoading" class="loading-screen">
          <div class="spinner"></div>
          <p>Metric Neo wird geladen...</p>
        </div>

        <!-- Setup-Dialog beim ersten Start -->
        <SetupDialog
          v-else-if="needsSetup"
          @setupComplete="onSetupComplete"
        />

        <!-- Haupt-Anwendung -->
        <MainApp v-else />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

.loading-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-screen p {
  font-size: 1.2em;
  opacity: 0.9;
}
</style>
