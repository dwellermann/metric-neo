<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { SelectDataDirectory, CompleteSetup } from '../../wailsjs/go/main/App';

const { t } = useI18n();
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
    
    if (result.success) {  // lowercase!
      console.log('Directory selected:', result.data);
      selectedDirectory.value = result.data;
    } else {
      console.error('Selection failed:', result.error);
      errorMessage.value = result.error || 'Fehler bei Verzeichnis-Auswahl';
    }
  } catch (error) {
    console.error('Exception caught:', error);
    errorMessage.value = 'Unerwarteter Fehler: ' + error.message;
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
    
    if (result.success) {  // lowercase!
      console.log('Setup completed successfully!');
      emit('setupComplete', selectedDirectory.value);
    } else {
      console.error('Setup failed:', result.error);
      errorMessage.value = result.error || 'Fehler beim Setup';
      // Zurück zur Verzeichnis-Auswahl
      selectedDirectory.value = '';
    }
  } catch (error) {
    console.error('Exception during setup:', error);
    errorMessage.value = 'Unerwarteter Fehler: ' + error.message;
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
    <div class="setup-dialog">
      <div class="setup-header">
        <img src="../assets/logo.svg" alt="Metric Neo Logo" class="logo" />
        <h1>🎯 {{ t('setup.welcome') }}</h1>
        <p class="subtitle">{{ t('setup.subtitle') }}</p>
      </div>

      <div class="setup-content">
        <!-- Schritt 1: Verzeichnis wählen -->
        <div v-if="!selectedDirectory" class="info-box">
          <p>
            <strong>{{ t('setup.firstSetup') }}</strong><br>
            {{ t('setup.selectDirectory') }}
          </p>
          <p class="note">
            💡 {{ t('setup.tipCreateFolder') }} <code>MetricNeo</code>)
          </p>
          <p class="note">
            {{ t('setup.dataStorage') }}
          </p>
        </div>

        <button 
          v-if="!selectedDirectory"
          @click="selectDirectory" 
          :disabled="isSelecting"
          class="select-button"
        >
          {{ isSelecting ? '⏳ ' + t('setup.buttonSelectingDir') : '📁 ' + t('setup.buttonSelectDir') }}
        </button>

        <!-- Schritt 2: Bestätigung & Start -->
        <div v-if="selectedDirectory" class="confirmation-box">
          <div class="selected-dir">
            <p class="label">{{ t('setup.selectedDirectory') }}</p>
            <p class="directory-path">📂 {{ selectedDirectory }}</p>
          </div>

          <div class="info-text">
            <p>{{ t('setup.foldersCreated') }}</p>
            <ul>
              <li>📁 profiles/</li>
              <li>📁 projectiles/</li>
              <li>📁 sessions/</li>
              <li>📁 sights/</li>
            </ul>
          </div>

          <div class="button-group">
            <button 
              @click="changeDirectory" 
              class="secondary-button"
              :disabled="isCompleting"
            >
              ← {{ t('setup.buttonChangeDir') }}
            </button>
            
            <button 
              @click="startApp" 
              :disabled="isCompleting"
              class="start-button"
            >
              {{ isCompleting ? '⏳ ' + t('setup.buttonStarting') : '🚀 ' + t('setup.buttonStart') }}
            </button>
          </div>
        </div>

        <div v-if="errorMessage" class="error-message">
          ⚠️ {{ errorMessage }}
        </div>
      </div>

      <div class="setup-footer">
        <small>{{ t('setup.configLocation') }} <code>~/.config/metric-neo/</code></small>
      </div>
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.setup-dialog {
  background: white;
  border-radius: 16px;
  max-width: 600px;
  width: 100%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.setup-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 40px 30px;
  text-align: center;
}

.setup-header .logo {
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
}

.setup-header h1 {
  margin: 0 0 10px 0;
  font-size: 2em;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  opacity: 0.9;
  font-size: 1.1em;
}

.setup-content {
  padding: 30px;
}

.info-box {
  background: #f8f9fa;
  border-left: 4px solid #667eea;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.info-box p {
  margin: 0 0 15px 0;
  line-height: 1.6;
}

.info-box p:last-child {
  margin-bottom: 0;
}

.suggestion {
  background: white;
  padding: 10px;
  border-radius: 4px;
  font-size: 0.95em;
}

.suggestion code {
  color: #667eea;
  font-weight: 600;
}

.note {
  font-size: 0.9em;
  color: #6c757d;
}

.select-button {
  width: 100%;
  padding: 16px 24px;
  font-size: 1.1em;
  font-weight: 600;
  color: white;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.select-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
}

.select-button:active:not(:disabled) {
  transform: translateY(0);
}

.select-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  margin-top: 20px;
  padding: 12px;
  background: #fee;
  color: #c33;
  border-radius: 4px;
  border-left: 4px solid #c33;
}

.confirmation-box {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  border: 2px solid #667eea;
}

.selected-dir {
  margin-bottom: 20px;
}

.selected-dir .label {
  margin: 0 0 8px 0;
  font-weight: 600;
  color: #495057;
}

.directory-path {
  margin: 0;
  padding: 12px;
  background: white;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.95em;
  color: #667eea;
  word-break: break-all;
}

.info-text {
  margin: 20px 0;
  padding: 15px;
  background: white;
  border-radius: 4px;
}

.info-text p {
  margin: 0 0 10px 0;
  font-weight: 600;
  color: #495057;
}

.info-text ul {
  margin: 0;
  padding-left: 20px;
  color: #6c757d;
}

.info-text li {
  margin: 5px 0;
}

.button-group {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

.secondary-button {
  flex: 1;
  padding: 12px 20px;
  font-size: 1em;
  font-weight: 600;
  color: #495057;
  background: white;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.secondary-button:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
}

.secondary-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.start-button {
  flex: 2;
  padding: 16px 24px;
  font-size: 1.2em;
  font-weight: 600;
  color: white;
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.start-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(40, 167, 69, 0.4);
}

.start-button:active:not(:disabled) {
  transform: translateY(0);
}

.start-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.setup-footer {
  background: #f8f9fa;
  padding: 15px 30px;
  text-align: center;
  color: #6c757d;
}

.setup-footer code {
  background: white;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.9em;
}
</style>
