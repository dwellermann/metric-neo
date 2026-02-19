<template>
  <div class="page">
    <n-card :segmented="true">
      <template #header>
        <div class="header">
          <span class="mdi mdi-chart-line"></span>
          <span>{{ t('sessions.title') || 'Sessions' }}</span>
        </div>
      </template>

      <template #header-extra>
        <n-button type="primary" @click="showCreateModal = true" :disabled="loadingSessions">
          <template #icon>
            <span class="mdi mdi-plus"></span>
          </template>
          {{ t('sessions.create') || 'New Session' }}
        </n-button>
      </template>

      <n-space align="center" style="margin-bottom: 12px;">
        <n-text>{{ t('sessions.chronoStatus') || 'Chrono' }}:</n-text>
        <n-tag :type="chronoConfigured ? 'success' : 'warning'" size="small">
          {{ chronoConfigured ? (t('sessions.chronoConfigured') || 'Configured') : (t('sessions.chronoNotConfigured') || 'Not configured') }}
        </n-tag>
        <n-text v-if="chronoConfigured" depth="3">
          {{ chronoConfig.port }} @ {{ chronoConfig.baudRate }}
        </n-text>
      </n-space>

      <!-- Filter Section -->
      <n-space align="center" style="margin-bottom: 12px;">
        <n-input type="date" v-model:value="filterStartDate" :placeholder="t('sessions.filterFrom') || 'From'" style="width: 160px;" />
        <n-input type="date" v-model:value="filterEndDate" :placeholder="t('sessions.filterTo') || 'To'" style="width: 160px;" />
        <n-text depth="3">{{ filteredSessions.length }} / {{ sessions.length }}</n-text>
        <n-button text type="primary" @click="filterStartDate = null; filterEndDate = null;" v-if="filterStartDate || filterEndDate">
          {{ t('common.clear') || 'Clear' }}
        </n-button>
      </n-space>

      <!-- Loading State -->
      <div v-if="loadingSessions" style="padding: 40px; text-align: center;">
        <n-spin size="large" />
      </div>

      <!-- Empty State -->
      <n-empty v-else-if="filteredSessions.length === 0" :description="filterStartDate || filterEndDate ? (t('sessions.noMatching') || 'No matching sessions') : (t('common.noData') || 'No sessions')" />

      <!-- Sessions Table -->
      <n-data-table v-else :columns="sessionColumns" :data="filteredSessions" :pagination="false" />
    </n-card>

    <!-- Create Session Modal -->
    <n-modal v-model:show="showCreateModal" preset="dialog" :title="t('sessions.create') || 'New Session'">
      <n-form ref="createFormRef" :model="createForm" :rules="createRules">
        <n-form-item :label="t('sessions.profile') || 'Profile'" path="profileId">
          <n-select v-model:value="createForm.profileId" :options="profileOptions" />
        </n-form-item>
        <n-form-item :label="t('sessions.projectile') || 'Projectile'" path="projectileId">
          <n-select v-model:value="createForm.projectileId" :options="projectileOptions" />
        </n-form-item>
        <n-form-item :label="t('sessions.temperature') || 'Temperature (°C)'" path="temperature">
          <n-input-number v-model:value="createForm.temperature" :step="0.1" />
        </n-form-item>
        <n-form-item :label="t('sessions.note') || 'Note'" path="note">
          <n-input v-model:value="createForm.note" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showCreateModal = false">{{ t('common.cancel') || 'Cancel' }}</n-button>
          <n-button type="primary" @click="handleCreateSession" :loading="creating">
            {{ t('common.create') || 'Create' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, h, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSelect,
  NSpace,
  NSpin,
  NTag,
  NText,
  useMessage,
  useDialog,
} from 'naive-ui';

const { t } = useI18n();
const router = useRouter();
const message = useMessage();
const dialog = useDialog();

const sessions = ref([]);
const profiles = ref([]);
const projectiles = ref([]);
const loadingSessions = ref(true);
const creating = ref(false);

const showCreateModal = ref(false);
const filterStartDate = ref(null);
const filterEndDate = ref(null);

const chronoConfig = ref({
  enabled: false,
  port: '',
  baudRate: 0,
  autoRecord: true,
});

const chronoConfigured = computed(() => {
  return chronoConfig.value.enabled && !!chronoConfig.value.port && chronoConfig.value.baudRate > 0;
});

const createFormRef = ref(null);
const createForm = ref({
  profileId: null,
  projectileId: null,
  temperature: null,
  note: '',
});

const createRules = {
  profileId: {
    required: true,
    message: t('validation.required') || 'Required',
    trigger: 'change',
  },
  projectileId: {
    required: true,
    message: t('validation.required') || 'Required',
    trigger: 'change',
  },
};

const profileOptions = computed(() =>
  profiles.value.map((p) => ({
    label: p.name,
    value: p.id,
  }))
);

const projectileOptions = computed(() =>
  projectiles.value.map((p) => ({
    label: p.name,
    value: p.id,
  }))
);

const formatDate = (iso) => {
  if (!iso) return '-';
  const date = new Date(iso);
  return date.toLocaleString();
};

const filteredSessions = computed(() => {
  if (!sessions.value || sessions.value.length === 0) return [];
  return sessions.value.filter((session) => {
    const sessionDate = new Date(session.createdAt);
    if (filterStartDate.value) {
      const start = new Date(filterStartDate.value);
      start.setHours(0, 0, 0, 0);
      if (sessionDate < start) return false;
    }
    if (filterEndDate.value) {
      const end = new Date(filterEndDate.value);
      end.setHours(23, 59, 59, 999);
      if (sessionDate > end) return false;
    }
    return true;
  });
});

const sessionColumns = [
  {
    title: t('sessions.createdAt') || 'Created',
    key: 'createdAt',
    render: (row) => formatDate(row.createdAt),
  },
  {
    title: t('sessions.profile') || 'Profile',
    key: 'profileName',
  },
  {
    title: t('sessions.projectile') || 'Projectile',
    key: 'projectileName',
  },
  {
    title: t('sessions.shots') || 'Shots',
    key: 'shotCount',
    render: (row) => `${row.validShotCount}/${row.shotCount}`,
  },
  {
    title: t('sessions.avgVelocity') || 'Avg Velocity',
    key: 'avgVelocityMPS',
    render: (row) => row.avgVelocityMPS ? `${row.avgVelocityMPS.toFixed(2)} m/s` : '-',
  },
  {
    title: t('sessions.avgEnergy') || 'Avg Energy',
    key: 'avgEnergyJoules',
    render: (row) => row.avgEnergyJoules ? `${row.avgEnergyJoules.toFixed(2)} J` : '-',
  },
  {
    title: t('sessions.actions') || 'Actions',
    key: 'actions',
    align: 'center',
    render: (row) => h(NSpace, null, {
      default: () => [
        h(
          NButton,
          {
            text: true,
            type: 'info',
            size: 'small',
            onClick: () => router.push(`/sessions/${row.id}`),
          },
          { default: () => h('span', { class: 'mdi mdi-eye' }) }
        ),
        h(
          NButton,
          {
            text: true,
            type: 'error',
            size: 'small',
            onClick: () => handleDeleteSession(row),
          },
          { default: () => h('span', { class: 'mdi mdi-delete-outline' }) }
        ),
      ],
    }),
  },
];

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

const loadChronoConfig = async () => {
  const fn = getBinding('GetChronoConfig');
  if (!fn) return;
  const result = await fn();
  const parsed = parseWailsResult(result);
  if (parsed?.success) {
    chronoConfig.value = parsed.data || chronoConfig.value;
  }
};

const loadProfiles = async () => {
  const fn = getBinding('ProfileListProfiles');
  if (!fn) return;
  const result = await fn();
  const parsed = parseWailsResult(result);
  if (parsed?.success) {
    profiles.value = parsed.data || [];
  }
};

const loadProjectiles = async () => {
  const fn = getBinding('ProjectileListProjectiles');
  if (!fn) return;
  const result = await fn();
  const parsed = parseWailsResult(result);
  if (parsed?.success) {
    projectiles.value = parsed.data || [];
  }
};

const loadSessions = async () => {
  loadingSessions.value = true;
  try {
    const fn = getBinding('SessionListSessions');
    if (!fn) {
      sessions.value = [];
      return;
    }
    const result = await fn();
    const parsed = parseWailsResult(result);
    if (parsed?.success) {
      sessions.value = parsed.data || [];
    } else {
      message.error(parsed?.error || t('common.error') || 'Error loading sessions');
    }
  } catch (err) {
    message.error(err.message || 'Failed to load sessions');
  } finally {
    loadingSessions.value = false;
  }
};

const handleCreateSession = async () => {
  if (!createFormRef.value) return;

  await createFormRef.value.validate();

  try {
    creating.value = true;
    const fn = getBinding('SessionCreateSession');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }

    const result = await fn(
      createForm.value.profileId,
      createForm.value.projectileId,
      createForm.value.temperature,
      createForm.value.note
    );

    if (result?.success) {
      const parsed = parseWailsResult(result);
      const newSessionId = parsed?.data?.id;

      message.success(t('common.created') || 'Session created');
      showCreateModal.value = false;
      createForm.value = {
        profileId: null,
        projectileId: null,
        temperature: null,
        note: '',
      };

      // Navigate to the new session directly
      if (newSessionId) {
        router.push(`/sessions/${newSessionId}`);
      } else {
        await loadSessions();
      }
    } else {
      message.error(result?.error);
    }
  } catch (err) {
    message.error(err.message || 'Failed to create session');
  } finally {
    creating.value = false;
  }
};

const handleDeleteSession = (row) => {
  dialog.warning({
    title: t('common.delete') || 'Delete',
    content: `Delete session "${row.profileName} · ${row.projectileName}"?`,
    positiveText: t('common.delete') || 'Delete',
    negativeText: t('common.cancel') || 'Cancel',
    onPositiveClick: () => deleteSession(row.id),
  });
};

const deleteSession = async (sessionId) => {
  try {
    const fn = getBinding('SessionDeleteSession');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }
    const result = await fn(sessionId);
    if (result?.success) {
      message.success(t('common.deleted') || 'Deleted');
      await loadSessions();
    } else {
      message.error(result?.error);
    }
  } catch (err) {
    message.error(err.message || 'Failed to delete session');
  }
};

onMounted(async () => {
  await new Promise(resolve => setTimeout(resolve, 100));
  await loadChronoConfig();
  await loadProfiles();
  await loadProjectiles();
  await loadSessions();
});
</script>

<style scoped>
.page {
  max-width: 1400px;
  margin: 0 auto;
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
</style>
