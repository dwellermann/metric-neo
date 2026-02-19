<template>
  <div class="page">
    <n-card :segmented="true">
      <template #header>
        <div class="header">
          <n-button text @click="$router.back()" size="small">
            <template #icon>
              <span class="mdi mdi-arrow-left"></span>
            </template>
          </n-button>
          <span class="mdi mdi-crosshairs"></span>
          <span>{{ t('sessions.detail') || 'Session Detail' }}</span>
        </div>
      </template>

      <n-space vertical size="large" style="max-width: 900px; margin: 0 auto;">
        <!-- Velocity Flash Display -->
        <transition name="flash">
          <div v-if="showFlash" class="velocity-flash">
            <div class="flash-content">
              <div class="flash-velocity">{{ formatNumber(flashVelocity, 2) }}</div>
              <div class="flash-unit">m/s</div>
            </div>
          </div>
        </transition>

        <!-- Header Info -->
        <div>
          <div class="detail-title">{{ session?.profileSnapshot?.name }} · {{ session?.projectileSnapshot?.name }}</div>
          <div class="detail-meta">
            {{ formatDate(session?.createdAt) }} · {{ session?.shots?.length || 0 }} {{ t('sessions.shots') || 'shots' }}
          </div>
        </div>

        <!-- Chrono Status -->
        <n-space align="center">
          <n-text>{{ t('sessions.chronoStatus') || 'Chrono' }}:</n-text>
          <n-tag :type="chronoConfigured ? 'success' : 'warning'" size="small">
            {{ chronoConfigured ? (t('sessions.chronoConfigured') || 'Configured') : (t('sessions.chronoNotConfigured') || 'Not configured') }}
          </n-tag>
          <n-text v-if="chronoConfigured" depth="3">
            {{ chronoConfig.port }} @ {{ chronoConfig.baudRate }}
          </n-text>
        </n-space>

        <!-- Recording Section -->
        <n-card size="small">
          <n-space vertical size="small" style="width: 100%;">
            <n-space align="center" justify="space-between">
              <n-text strong>{{ t('sessions.recordShot') || 'Record Shot' }}</n-text>
              <n-tag v-if="chronoListening" size="small" type="success" style="animation: pulse 1s infinite;">
                <span class="mdi mdi-pulse"></span>
                {{ t('sessions.listening') || 'Listening...' }}
              </n-tag>
            </n-space>

            <n-space align="center">
              <n-input-number v-model:value="recordVelocity" :min="0" :step="0.1" style="width: 160px;" />
              <n-text>m/s</n-text>
            </n-space>

            <n-space>
              <n-button
                type="primary"
                @click="handleRecordShot"
                :loading="recording"
                :disabled="!session"
              >
                {{ t('sessions.record') || 'Record' }}
              </n-button>
              <n-button
                v-if="!chronoListening"
                type="info"
                @click="startMeasurement"
                :disabled="!chronoConfigured"
              >
                {{ t('sessions.startMeasurement') || 'Start Measurement' }}
              </n-button>
              <n-button
                v-if="chronoListening"
                type="warning"
                @click="stopMeasurement"
              >
                {{ t('sessions.stopMeasurement') || 'Stop' }}
              </n-button>
            </n-space>

            <n-text depth="3" style="display: block; margin-top: 4px;">
              {{ chronoConfigured ? (t('sessions.chronoReady') || 'Chrono ready') : (t('sessions.chronoManual') || 'Manual input') }}
            </n-text>
          </n-space>
        </n-card>

        <!-- Statistics -->
        <n-card size="small">
          <n-text strong>{{ t('sessions.statistics') || 'Statistics' }}</n-text>
          <n-grid :cols="3" :x-gap="12" :y-gap="12" style="margin-top: 12px;">
            <n-gi>
              <div class="stat-item">
                <div class="stat-label">{{ t('sessions.avgVelocity') || 'Avg Velocity' }}</div>
                <div class="stat-value">{{ formatNumber(stats.avgVelocityMPS, 2) }} m/s</div>
              </div>
            </n-gi>
            <n-gi>
              <div class="stat-item">
                <div class="stat-label">{{ t('sessions.stdDev') || 'Std. Deviation' }}</div>
                <div class="stat-value">{{ formatNumber(stats.standardDeviation, 2) }} m/s</div>
              </div>
            </n-gi>
            <n-gi>
              <div class="stat-item">
                <div class="stat-label">{{ t('sessions.extremeSpread') || 'Extreme Spread' }}</div>
                <div class="stat-value">{{ formatNumber(stats.extremeSpread, 2) }} m/s</div>
              </div>
            </n-gi>
            <n-gi>
              <div class="stat-item">
                <div class="stat-label">{{ t('sessions.minVelocity') || 'Min' }}</div>
                <div class="stat-value">{{ formatNumber(stats.minVelocityMPS, 2) }} m/s</div>
              </div>
            </n-gi>
            <n-gi>
              <div class="stat-item">
                <div class="stat-label">{{ t('sessions.maxVelocity') || 'Max' }}</div>
                <div class="stat-value">{{ formatNumber(stats.maxVelocityMPS, 2) }} m/s</div>
              </div>
            </n-gi>
            <n-gi>
              <div class="stat-item">
                <div class="stat-label">{{ t('sessions.avgEnergy') || 'Avg Energy' }}</div>
                <div class="stat-value">{{ formatNumber(stats.avgEnergyJoules, 2) }} J</div>
              </div>
            </n-gi>
          </n-grid>
        </n-card>

        <!-- Shots Table -->
        <n-card size="small">
          <n-text strong>{{ t('sessions.shots') || 'Shots' }} ({{ session?.shots?.length || 0 }})</n-text>
          <div v-if="!session?.shots?.length" style="margin-top: 12px; padding: 20px; text-align: center;">
            <n-empty :description="t('common.noData') || 'No shots'" />
          </div>
          <n-data-table
            v-else
            :columns="shotColumns"
            :data="session.shots"
            :pagination="false"
            style="margin-top: 12px;"
          />
        </n-card>

        <!-- Note -->
        <n-card size="small">
          <n-text strong>{{ t('sessions.note') || 'Note' }}</n-text>
          <n-space align="center" style="margin-top: 8px; width: 100%;">
            <n-input
              v-model:value="noteEdit"
              type="textarea"
              :autosize="{ minRows: 2, maxRows: 4 }"
              style="flex: 1;"
            />
            <n-button type="primary" @click="handleSaveNote" :loading="savingNote">
              {{ t('common.save') || 'Save' }}
            </n-button>
          </n-space>
        </n-card>

        <!-- Delete Button -->
        <n-space>
          <n-button type="error" @click="handleDelete">
            {{ t('common.delete') || 'Delete' }}
          </n-button>
        </n-space>
      </n-space>
    </n-card>
  </div>
</template>

<script setup>
import { ref, h, onMounted, onBeforeUnmount, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NGrid,
  NGi,
  NInput,
  NInputNumber,
  NSpace,
  NTag,
  NText,
  useDialog,
  useMessage,
} from 'naive-ui';

const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const dialog = useDialog();
const message = useMessage();

const session = ref(null);
const stats = ref({
  avgVelocityMPS: 0,
  standardDeviation: 0,
  minVelocityMPS: 0,
  maxVelocityMPS: 0,
  extremeSpread: 0,
  avgEnergyJoules: 0,
  validShotCount: 0,
  totalShotCount: 0,
});

const chronoConfig = ref({
  enabled: false,
  port: '',
  baudRate: 0,
  autoRecord: true,
});

const chronoConfigured = computed(() => {
  return chronoConfig.value.enabled && !!chronoConfig.value.port && chronoConfig.value.baudRate > 0;
});

const recordVelocity = ref(null);
const noteEdit = ref('');
const recording = ref(false);
const savingNote = ref(false);
const chronoListening = ref(false);
let chronoPollTimer = null;

// Flash display for shot velocity
const flashVelocity = ref(null);
const showFlash = ref(false);
let flashTimer = null;

const getBinding = (method) => {
  if (!window.go?.main?.App) return null;
  return window.go.main.App[method];
};

const parseWailsResult = (result) => {
  if (!result?.data) return result;

  if (Array.isArray(result.data)) {
    return {
      ...result,
      data: result.data.map((item) => {
        const stringified = JSON.stringify(item);
        if (stringified.includes('\\\"') || (stringified.startsWith('"{') && stringified.endsWith('}"'))) {
          try {
            const firstParse = typeof item === 'string' ? item : JSON.parse(item);
            const secondParse = JSON.parse(firstParse);
            return secondParse;
          } catch (e) {
            return item;
          }
        }

        if (typeof item === 'string') {
          try {
            return JSON.parse(item);
          } catch (e) {
            return item;
          }
        }

        return item;
      }),
    };
  } else if (typeof result.data === 'string') {
    try {
      return {
        ...result,
        data: JSON.parse(result.data)
      };
    } catch (e) {
      return result;
    }
  }
  return result;
};

const showVelocityFlash = (velocity) => {
  if (flashTimer) {
    clearTimeout(flashTimer);
  }
  flashVelocity.value = velocity;
  showFlash.value = true;
  flashTimer = setTimeout(() => {
    showFlash.value = false;
    flashVelocity.value = null;
  }, 5000);
};

const formatNumber = (value, decimals = 2) => {
  if (value === null || value === undefined || Number.isNaN(value)) return '-';
  return Number(value).toFixed(decimals);
};

const formatDate = (iso) => {
  if (!iso) return '-';
  const date = new Date(iso);
  return date.toLocaleString();
};

const shotColumns = [
  {
    title: t('common.add') || '#',
    key: 'index',
    render: (row, index) => index + 1,
    width: 50,
  },
  {
    title: t('sessions.velocity') || 'Velocity',
    key: 'velocityMPS',
    render: (row) => formatNumber(row.velocityMPS, 2) + ' m/s',
  },
  {
    title: t('sessions.energy') || 'Energy',
    key: 'energyJoules',
    render: (row) => formatNumber(row.energyJoules, 2) + ' J',
  },
  {
    title: t('sessions.timestamp') || 'Timestamp',
    key: 'timestamp',
    render: (row) => formatDate(row.timestamp),
  },
  {
    title: t('sessions.valid') || 'Valid',
    key: 'valid',
    render: (row) => row.valid ? '✓' : '✗',
  },
  {
    title: t('sessions.actions') || 'Actions',
    key: 'actions',
    align: 'center',
    render: (row, index) => h(NSpace, null, {
      default: () => [
        h(
          NButton,
          {
            text: true,
            type: 'warning',
            size: 'small',
            disabled: !row.valid,
            onClick: () => handleMarkInvalid(index),
          },
          { default: () => h('span', { class: 'mdi mdi-close-circle-outline' }) }
        ),
      ],
    }),
  },
];

const loadChronoConfig = async () => {
  const fn = getBinding('GetChronoConfig');
  if (!fn) return;
  const result = await fn();
  const parsed = parseWailsResult(result);
  if (parsed?.success) {
    chronoConfig.value = parsed.data || chronoConfig.value;
  }
};

const loadSessionDetail = async (id) => {
  const fn = getBinding('SessionLoadSession');
  if (!fn) return;
  const result = await fn(id);
  const parsed = parseWailsResult(result);
  if (parsed?.success) {
    session.value = parsed.data;
    noteEdit.value = parsed.data.note || '';
  } else {
    message.error(parsed?.error || t('common.error') || 'Error loading session');
  }
};

const loadStatistics = async (id) => {
  const fn = getBinding('SessionGetStatistics');
  if (!fn) return;
  const result = await fn(id);
  const parsed = parseWailsResult(result);
  if (parsed?.success) {
    stats.value = parsed.data || stats.value;
  }
};

const pollChronoOnce = async () => {
  if (!session.value) return;
  const fn = getBinding('SessionPollChrono');
  if (!fn) return;
  try {
    const result = await fn(session.value.id);
    const parsed = parseWailsResult(result);
    if (!parsed?.success) {
      if (parsed?.error) {
        message.error(parsed.error);
      }
      return;
    }

    if (parsed.data?.recorded && parsed.data?.session) {
      session.value = parsed.data.session;
      noteEdit.value = parsed.data.session.note || '';
      await loadStatistics(session.value.id);

      // Show flash for newly recorded shot
      if (parsed.data?.velocityMPS) {
        showVelocityFlash(parsed.data.velocityMPS);
      }
    }
  } catch (err) {
    message.error(err.message || 'Failed to poll chrono');
  }
};

const startMeasurement = () => {
  if (chronoListening.value) return;
  chronoListening.value = true;
  chronoPollTimer = setInterval(pollChronoOnce, 500);
};

const stopMeasurement = () => {
  if (chronoPollTimer) {
    clearInterval(chronoPollTimer);
    chronoPollTimer = null;
  }
  chronoListening.value = false;
};

const handleRecordShot = async () => {
  if (!session.value) return;
  if (recordVelocity.value === null || recordVelocity.value === undefined) {
    message.error(t('validation.required') || 'Required');
    return;
  }
  try {
    recording.value = true;
    const fn = getBinding('SessionRecordShot');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }
    const result = await fn(session.value.id, recordVelocity.value);
    if (result?.success) {
      showVelocityFlash(recordVelocity.value);
      recordVelocity.value = null;
      await loadSessionDetail(session.value.id);
      await loadStatistics(session.value.id);
    } else {
      message.error(result?.error);
    }
  } catch (err) {
    message.error(err.message || 'Failed to record shot');
  } finally {
    recording.value = false;
  }
};

const handleMarkInvalid = async (index) => {
  if (!session.value) return;
  const fn = getBinding('SessionMarkShotInvalid');
  if (!fn) {
    message.error('Backend not ready');
    return;
  }
  const result = await fn(session.value.id, index);
  if (result?.success) {
    await loadSessionDetail(session.value.id);
    await loadStatistics(session.value.id);
  } else {
    message.error(result?.error);
  }
};

const handleSaveNote = async () => {
  if (!session.value) return;
  try {
    savingNote.value = true;
    const fn = getBinding('SessionUpdateNote');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }
    const result = await fn(session.value.id, noteEdit.value || '');
    if (result?.success) {
      await loadSessionDetail(session.value.id);
    } else {
      message.error(result?.error);
    }
  } catch (err) {
    message.error(err.message || 'Failed to save note');
  } finally {
    savingNote.value = false;
  }
};

const handleDelete = () => {
  dialog.warning({
    title: t('common.delete') || 'Delete',
    content: `${t('common.deleteConfirm') || 'Delete'} Session?`,
    positiveText: t('common.delete') || 'Delete',
    negativeText: t('common.cancel') || 'Cancel',
    onPositiveClick: async () => {
      const fn = getBinding('SessionDeleteSession');
      if (!fn) {
        message.error('Backend not ready');
        return;
      }
      const result = await fn(session.value.id);
      if (result?.success) {
        router.back();
      } else {
        message.error(result?.error);
      }
    },
  });
};

onMounted(async () => {
  await new Promise(resolve => setTimeout(resolve, 100));
  const sessionId = route.params.id;
  if (sessionId) {
    await loadChronoConfig();
    await loadSessionDetail(sessionId);
    await loadStatistics(sessionId);
  }
});

onBeforeUnmount(() => {
  stopMeasurement();
  if (flashTimer) {
    clearTimeout(flashTimer);
  }
});
</script>

<style scoped>
.page {
  max-width: 100%;
}

.header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 600;
}

.header .mdi {
  font-size: 1.2rem;
}

.detail-title {
  font-weight: 600;
  font-size: 1.1rem;
}

.detail-meta {
  font-size: 0.9rem;
  color: var(--n-text-color-3);
  margin-top: 4px;
}

.stat-item {
  padding: 12px;
  border: 1px solid var(--n-border-color);
  border-radius: 4px;
  text-align: center;
}

.stat-label {
  font-size: 0.85rem;
  color: var(--n-text-color-3);
  margin-bottom: 4px;
}

.stat-value {
  font-weight: 600;
  font-size: 1.1rem;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

.velocity-flash {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 9999;
  background: linear-gradient(135deg, #18a058 0%, #36ad6a 100%);
  border-radius: 24px;
  padding: 48px 64px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  border: 4px solid rgba(255, 255, 255, 0.2);
}

.flash-content {
  text-align: center;
}

.flash-velocity {
  font-size: 6rem;
  font-weight: 700;
  color: white;
  line-height: 1;
  text-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.flash-unit {
  font-size: 2rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  margin-top: 8px;
}

.flash-enter-active {
  animation: flashIn 0.3s ease-out;
}

.flash-leave-active {
  animation: flashOut 0.3s ease-in;
}

@keyframes flashIn {
  from {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.8);
  }
  to {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
}

@keyframes flashOut {
  from {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
  to {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.8);
  }
}
</style>
