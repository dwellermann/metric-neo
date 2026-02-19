<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { SelectDataDirectory, CompleteSetup } from '../../wailsjs/go/main/App';
import {
  NButton,
  NCard,
  NSpace,
  NText,
  NAlert,
  NCode,
  useMessage,
} from 'naive-ui';

const { t } = useI18n();
const message = useMessage();
const emit = defineEmits(['setupComplete']);

const isSelecting = ref(false);
const selectedDirectory = ref('');
const errorMessage = ref('');
const isCompleting = ref(false);

async function selectDirectory() {
  console.log('Starting directory selection...');
  isSelecting.value = true;
  errorMessage.value = '';

  try {
    console.log('Calling SelectDataDirectory()...');
    const result = await SelectDataDirectory();
    console.log('Result received:', result);

    if (result.success) {
      console.log('Directory selected:', result.data);
      selectedDirectory.value = result.data;
      message.success(t('setup.directorySelected') || 'Directory selected');
    } else {
      console.error('Selection failed:', result.error);
      errorMessage.value = result.error || 'Fehler bei Verzeichnis-Auswahl';
      message.error(errorMessage.value);
    }
  } catch (error) {
    console.error('Exception caught:', error);
    errorMessage.value = 'Unerwarteter Fehler: ' + error.message;
    message.error(errorMessage.value);
  } finally {
    isSelecting.value = false;
  }
}

async function startApp() {
  console.log('Starting app with directory:', selectedDirectory.value);
  isCompleting.value = true;
  errorMessage.value = '';

  try {
    const result = await CompleteSetup(selectedDirectory.value);

    if (result.success) {
      console.log('Setup completed successfully!');
      message.success(t('setup.setupComplete') || 'Setup complete!');
      emit('setupComplete', selectedDirectory.value);
    } else {
      console.error('Setup failed:', result.error);
      errorMessage.value = result.error || 'Fehler beim Setup';
      message.error(errorMessage.value);
      selectedDirectory.value = '';
    }
  } catch (error) {
    console.error('Exception during setup:', error);
    errorMessage.value = 'Unerwarteter Fehler: ' + error.message;
    message.error(errorMessage.value);
    selectedDirectory.value = '';
  } finally {
    isCompleting.value = false;
  }
}

function changeDirectory() {
  selectedDirectory.value = '';
  errorMessage.value = '';
}
</script>

<template>
  <div class="setup-overlay">
    <div class="setup-container">
      <n-card size="large" class="setup-card">
        <n-space vertical :size="24" align="center">
          <!-- Logo & Header -->
          <div class="header-section">
            <img src="../assets/logo.svg" alt="Metric Neo Logo" class="logo" />
            <h1 class="title">{{ t('setup.welcome') }}</h1>
            <p class="subtitle">{{ t('setup.subtitle') }}</p>
          </div>

          <!-- Content -->
          <div class="content-section">
            <!-- Schritt 1: Verzeichnis wählen -->
            <div v-if="!selectedDirectory">
              <n-alert type="info" style="margin-bottom: 20px;">
                <n-space vertical :size="8">
                  <n-text strong>{{ t('setup.firstSetup') }}</n-text>
                  <n-text>{{ t('setup.selectDirectory') }}</n-text>
                  <n-text depth="3" style="font-size: 0.9em;">
                    💡 {{ t('setup.tipCreateFolder') }} <n-code>MetricNeo</n-code>)
                  </n-text>
                  <n-text depth="3" style="font-size: 0.9em;">
                    {{ t('setup.dataStorage') }}
                  </n-text>
                </n-space>
              </n-alert>

              <n-button
                type="primary"
                size="large"
                block
                @click="selectDirectory"
                :loading="isSelecting"
              >
                <template #icon>
                  <span class="mdi mdi-folder-open"></span>
                </template>
                {{ t('setup.buttonSelectDir') }}
              </n-button>
            </div>

            <!-- Schritt 2: Bestätigung & Start -->
            <div v-if="selectedDirectory">
              <n-alert type="success" style="margin-bottom: 16px;">
                <n-space vertical :size="8">
                  <n-text strong>{{ t('setup.selectedDirectory') }}</n-text>
                  <n-code style="word-break: break-all; display: block;">
                    {{ selectedDirectory }}
                  </n-code>
                </n-space>
              </n-alert>

              <n-card size="small" embedded style="margin-bottom: 16px;">
                <n-text strong style="display: block; margin-bottom: 12px;">
                  {{ t('setup.foldersCreated') }}
                </n-text>
                <n-space vertical :size="4">
                  <n-text depth="3">
                    <span class="mdi mdi-folder" style="margin-right: 8px;"></span>
                    profiles/
                  </n-text>
                  <n-text depth="3">
                    <span class="mdi mdi-folder" style="margin-right: 8px;"></span>
                    projectiles/
                  </n-text>
                  <n-text depth="3">
                    <span class="mdi mdi-folder" style="margin-right: 8px;"></span>
                    sessions/
                  </n-text>
                  <n-text depth="3">
                    <span class="mdi mdi-folder" style="margin-right: 8px;"></span>
                    sights/
                  </n-text>
                </n-space>
              </n-card>

              <n-space justify="space-between">
                <n-button
                  @click="changeDirectory"
                  :disabled="isCompleting"
                  secondary
                >
                  <template #icon>
                    <span class="mdi mdi-arrow-left"></span>
                  </template>
                  {{ t('setup.buttonChangeDir') }}
                </n-button>

                <n-button
                  type="success"
                  size="large"
                  @click="startApp"
                  :loading="isCompleting"
                >
                  <template #icon>
                    <span class="mdi mdi-rocket-launch"></span>
                  </template>
                  {{ t('setup.buttonStart') }}
                </n-button>
              </n-space>
            </div>

            <!-- Error -->
            <n-alert v-if="errorMessage" type="error" style="margin-top: 16px;">
              {{ errorMessage }}
            </n-alert>
          </div>

          <!-- Footer -->
          <div class="footer-section">
            <n-text depth="3" style="text-align: center; font-size: 0.85em;">
              {{ t('setup.configLocation') }} <n-code>~/.config/metric-neo/</n-code>
            </n-text>
          </div>
        </n-space>
      </n-card>
    </div>
  </div>
</template>

<style scoped>
.setup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  z-index: 10000;
  overflow-y: auto;
  overflow-x: hidden;
}

.setup-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.setup-card {
  max-width: 600px;
  width: 100%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  border-radius: 16px;
}

.header-section {
  text-align: center;
  width: 100%;
  padding-bottom: 10px;
}

.logo {
  width: 80px;
  height: 80px;
  display: block;
  margin: 0 auto 20px auto;
}

.title {
  margin: 0 0 8px 0;
  font-size: 1.8em;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  font-size: 1em;
  opacity: 0.7;
}

.content-section {
  width: 100%;
}

.footer-section {
  width: 100%;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}
</style>
