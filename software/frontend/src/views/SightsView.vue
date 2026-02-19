<template>
  <div class="page">
    <n-card :segmented="true">
      <template #header>
        <div class="header">
          <span class="mdi mdi-target"></span>
          <span>{{ t('sights.title') || 'Sights' }}</span>
        </div>
      </template>

      <template #header-extra>
        <n-button type="primary" @click="showCreateModal = true" :disabled="loading">
          <template #icon>
            <span class="mdi mdi-plus"></span>
          </template>
          {{ t('common.add') || 'Add' }}
        </n-button>
      </template>

      <!-- Loading State -->
      <div v-if="loading" style="padding: 40px; text-align: center;">
        <n-spin size="large" />
      </div>

      <!-- Empty State -->
      <n-empty v-else-if="sights.length === 0" :description="t('common.noData') || 'No sights'" />

      <!-- Sights Table -->
      <n-data-table v-else :columns="columns" :data="sights" :pagination="false" />
    </n-card>

    <!-- Create Sight Modal -->
    <n-modal v-model:show="showCreateModal" preset="dialog" :title="t('sights.createNew') || 'New Sight'">
      <n-form ref="createFormRef" :model="createForm" :rules="formRules">
        <n-form-item :label="t('sights.type') || 'Type'" path="type">
          <n-select v-model:value="createForm.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item :label="t('sights.modelName') || 'Model Name'" path="modelName">
          <n-input v-model:value="createForm.modelName" placeholder="e.g., Leupold Mark 5HD" />
        </n-form-item>
        <n-form-item :label="t('sights.weight') || 'Weight (g)'" path="weight">
          <n-input-number v-model:value="createForm.weight" :min="0" :step="1" />
        </n-form-item>
        <n-form-item :label="t('sights.minMagnification') || 'Min Magnification'" path="minMagnification">
          <n-input-number v-model:value="createForm.minMagnification" :min="0" :step="0.1" />
        </n-form-item>
        <n-form-item :label="t('sights.maxMagnification') || 'Max Magnification'" path="maxMagnification">
          <n-input-number v-model:value="createForm.maxMagnification" :min="0" :step="0.1" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showCreateModal = false">{{ t('common.cancel') || 'Cancel' }}</n-button>
          <n-button type="primary" @click="handleCreateSight" :loading="creating">
            {{ t('common.create') || 'Create' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- Edit Sight Modal -->
    <n-modal v-model:show="showEditModal" preset="dialog" :title="t('sights.edit') || 'Edit Sight'">
      <n-form v-if="editingSight" ref="editFormRef" :model="editForm" :rules="formRules">
        <n-form-item :label="t('sights.type') || 'Type'" path="type">
          <n-select v-model:value="editForm.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item :label="t('sights.modelName') || 'Model Name'" path="modelName">
          <n-input v-model:value="editForm.modelName" />
        </n-form-item>
        <n-form-item :label="t('sights.weight') || 'Weight (g)'" path="weight">
          <n-input-number v-model:value="editForm.weight" :min="0" :step="1" />
        </n-form-item>
        <n-form-item :label="t('sights.minMagnification') || 'Min Magnification'" path="minMagnification">
          <n-input-number v-model:value="editForm.minMagnification" :min="0" :step="0.1" />
        </n-form-item>
        <n-form-item :label="t('sights.maxMagnification') || 'Max Magnification'" path="maxMagnification">
          <n-input-number v-model:value="editForm.maxMagnification" :min="0" :step="0.1" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showEditModal = false">{{ t('common.cancel') || 'Cancel' }}</n-button>
          <n-button type="primary" @click="handleSaveSight" :loading="saving">
            {{ t('common.save') || 'Save' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, h, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  NButton,
  NCard,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSelect,
  NDataTable,
  NSpin,
  NSpace,
  useDialog,
  useMessage,
} from 'naive-ui';

const { t } = useI18n();
const dialog = useDialog();
const message = useMessage();

// Helper to get Wails bindings - they're injected at runtime
const getBinding = (method) => {
  if (!window.go?.main?.App) {
    return null;
  }
  return window.go.main.App[method];
};

// Helper to parse Wails Result data (handles double-stringification)
const parseWailsResult = (result) => {
  if (!result?.data) return result;

  if (Array.isArray(result.data)) {
    const parsed = {
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
    return parsed;
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

const sights = ref([]);
const loading = ref(true);
const creating = ref(false);
const saving = ref(false);

const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingSight = ref(null);

const createFormRef = ref(null);
const editFormRef = ref(null);

const typeOptions = [
  { label: 'Scope', value: 'scope' },
  { label: 'Red Dot', value: 'red_dot' },
  { label: 'Diopter', value: 'diopter' },
  { label: 'Open Sights', value: 'open_sights' },
];

const createForm = ref({
  type: 'scope',
  modelName: '',
  weight: 0,
  minMagnification: 1,
  maxMagnification: 10,
});

const editForm = ref({
  type: 'scope',
  modelName: '',
  weight: 0,
  minMagnification: 1,
  maxMagnification: 10,
});

const formRules = {
  modelName: {
    required: true,
    message: t('validation.required') || 'Model name is required',
    trigger: 'blur',
  },
  type: {
    required: true,
    message: t('validation.required') || 'Type is required',
    trigger: 'change',
  },
};

const columns = [
  {
    title: t('sights.type') || 'Type',
    key: 'type',
    render: (row) => {
      const typeMap = {
        scope: 'Scope',
        red_dot: 'Red Dot',
        diopter: 'Diopter',
        open_sights: 'Open Sights',
      };
      return typeMap[row.type] || row.type;
    },
  },
  {
    title: t('sights.modelName') || 'Model',
    key: 'modelName',
  },
  {
    title: t('sights.weight') || 'Weight',
    key: 'weightG',
    render: (row) => row.weightG ? `${row.weightG.toFixed(0)} g` : '-',
  },
  {
    title: t('sights.magnification') || 'Magnification',
    key: 'magnification',
    render: (row) => {
      if (row.minMagnification && row.maxMagnification) {
        return `${row.minMagnification.toFixed(1)}x - ${row.maxMagnification.toFixed(1)}x`;
      }
      return '-';
    },
  },
  {
    title: 'Actions',
    key: 'actions',
    align: 'center',
    render: (row) => h(NSpace, null, {
      default: () => [
        h(
          NButton,
          {
            text: true,
            type: 'primary',
            size: 'small',
            onClick: () => openEditModal(row),
          },
          { default: () => h('span', { class: 'mdi mdi-pencil' }) }
        ),
        h(
          NButton,
          {
            text: true,
            type: 'error',
            size: 'small',
            onClick: () => handleDeleteSight(row),
          },
          { default: () => h('span', { class: 'mdi mdi-delete' }) }
        ),
      ],
    }),
  },
];

// Load sights
const loadSights = async () => {
  loading.value = true;
  try {
    const fn = getBinding('SightListSights');
    if (!fn) {
      sights.value = [];
      return;
    }
    const result = await fn();
    const parsed = parseWailsResult(result);
    if (parsed?.success) {
      sights.value = parsed.data || [];
    } else {
      message.error(parsed?.error || t('common.error') || 'Error loading sights');
    }
  } catch (err) {
    message.error(err.message || 'Failed to load sights');
  } finally {
    loading.value = false;
  }
};

// Create sight
const handleCreateSight = async () => {
  if (!createFormRef.value) return;

  try {
    await createFormRef.value.validate();
    creating.value = true;

    const fn = getBinding('SightCreateSight');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }

    const result = await fn(
      createForm.value.type,
      createForm.value.modelName,
      createForm.value.weight,
      createForm.value.minMagnification,
      createForm.value.maxMagnification
    );

    if (result?.success) {
      message.success(t('common.created') || 'Sight created');
      showCreateModal.value = false;
      createForm.value = {
        type: 'scope',
        modelName: '',
        weight: 0,
        minMagnification: 1,
        maxMagnification: 10,
      };
      await loadSights();
    } else {
      message.error(result?.error);
    }
  } catch (err) {
    message.error(err.message);
  } finally {
    creating.value = false;
  }
};

// Open edit modal
const openEditModal = (sight) => {
  editingSight.value = sight;
  editForm.value = {
    type: sight.type || 'scope',
    modelName: sight.modelName || '',
    weight: sight.weightG || 0,
    minMagnification: sight.minMagnification || 1,
    maxMagnification: sight.maxMagnification || 10,
  };
  showEditModal.value = true;
};

// Save sight
const handleSaveSight = async () => {
  if (!editingSight.value) return;

  try {
    if (editFormRef.value) {
      await editFormRef.value.validate();
    }
    saving.value = true;

    const fn = getBinding('SightUpdateSight');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }

    const result = await fn(
      editingSight.value.id,
      editForm.value.type,
      editForm.value.modelName,
      editForm.value.weight,
      editForm.value.minMagnification,
      editForm.value.maxMagnification
    );

    if (result?.success) {
      message.success(t('common.saved') || 'Sight updated');
      showEditModal.value = false;
      await loadSights();
    } else {
      message.error(result?.error);
    }
  } catch (err) {
    message.error(err.message);
  } finally {
    saving.value = false;
  }
};

// Delete sight
const handleDeleteSight = (sight) => {
  dialog.warning({
    title: t('common.delete') || 'Delete',
    content: `${t('common.deleteConfirm') || 'Delete'} "${sight.modelName}"?`,
    positiveText: t('common.delete') || 'Delete',
    negativeText: t('common.cancel') || 'Cancel',
    onPositiveClick: async () => {
      try {
        const fn = getBinding('SightDeleteSight');
        if (!fn) {
          message.error('Backend not ready');
          return;
        }
        const result = await fn(sight.id);
        if (result?.success) {
          message.success(t('common.deleted') || 'Deleted');
          await loadSights();
        } else {
          message.error(result?.error);
        }
      } catch (err) {
        message.error(err.message);
      }
    },
  });
};

// Load on mount
onMounted(async () => {
  // Wait a bit for Wails to inject window.go
  await new Promise(resolve => setTimeout(resolve, 100));
  await loadSights();
});
</script>

<style scoped>
.page {
  max-width: 1200px;
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
