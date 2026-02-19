<template>
  <div class="page">
    <n-card :segmented="true">
      <template #header>
        <div class="header">
          <span class="mdi mdi-ammunition"></span>
          <span>{{ t('projectiles.title') }}</span>
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
      <n-empty v-else-if="projectiles.length === 0" :description="t('common.noData') || 'No projectiles'" />

      <!-- Projectiles Table -->
      <n-data-table v-else :columns="columns" :data="projectiles" :pagination="false" />
    </n-card>

    <!-- Create Projectile Modal -->
    <n-modal v-model:show="showCreateModal" preset="dialog" :title="t('projectiles.createNew') || 'New Projectile'">
      <n-form ref="createFormRef" :model="createForm" :rules="formRules">
        <n-form-item :label="t('projectiles.name') || 'Name'" path="name">
          <n-input v-model:value="createForm.name" placeholder="e.g., .308 Match" />
        </n-form-item>
        <n-form-item :label="t('projectiles.weight') || 'Weight (g)'" path="weight">
          <n-input-number v-model:value="createForm.weight" :min="0" :step="0.001" :precision="3" />
        </n-form-item>
        <n-form-item :label="t('projectiles.bc') || 'BC'" path="bc">
          <n-input-number v-model:value="createForm.bc" :min="0" :max="1" :step="0.001" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showCreateModal = false">{{ t('common.cancel') || 'Cancel' }}</n-button>
          <n-button type="primary" @click="handleCreateProjectile" :loading="creating">
            {{ t('common.create') || 'Create' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- Edit Projectile Modal -->
    <n-modal v-model:show="showEditModal" preset="dialog" :title="t('projectiles.edit') || 'Edit Projectile'">
      <n-form v-if="editingProjectile" ref="editFormRef" :model="editForm" :rules="formRules">
        <n-form-item :label="t('projectiles.name') || 'Name'" path="name">
          <n-input v-model:value="editForm.name" />
        </n-form-item>
        <n-form-item :label="t('projectiles.weight') || 'Weight (g)'" path="weight">
          <n-input-number v-model:value="editForm.weight" :min="0" :step="0.001" :precision="3" />
        </n-form-item>
        <n-form-item :label="t('projectiles.bc') || 'BC'" path="bc">
          <n-input-number v-model:value="editForm.bc" :min="0" :max="1" :step="0.001" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showEditModal = false">{{ t('common.cancel') || 'Cancel' }}</n-button>
          <n-button type="primary" @click="handleSaveProjectile" :loading="saving">
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

const projectiles = ref([]);
const loading = ref(true);
const creating = ref(false);
const saving = ref(false);

const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingProjectile = ref(null);

const createFormRef = ref(null);
const editFormRef = ref(null);

const createForm = ref({
  name: '',
  weight: 0,
  bc: 0.5,
});

const editForm = ref({
  name: '',
  weight: 0,
  bc: 0.5,
});

const formRules = {
  name: {
    required: true,
    message: t('validation.required') || 'Name is required',
    trigger: 'blur',
  },
};

const columns = [
  {
    title: t('projectiles.name') || 'Name',
    key: 'name',
  },
  {
    title: t('projectiles.weight') || 'Weight',
    key: 'weightGrams',
    render: (row) => row.weightGrams ? `${row.weightGrams.toFixed(3)} g` : '-',
  },
  {
    title: t('projectiles.bc') || 'BC',
    key: 'bc',
    render: (row) => row.bc ? `${row.bc.toFixed(3)}` : '-',
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
            onClick: () => handleDeleteProjectile(row),
          },
          { default: () => h('span', { class: 'mdi mdi-delete' }) }
        ),
      ],
    }),
  },
];

// Load projectiles
const loadProjectiles = async () => {
  loading.value = true;
  try {
    const fn = getBinding('ProjectileListProjectiles');
    if (!fn) {
      projectiles.value = [];
      return;
    }
    const result = await fn();
    const parsed = parseWailsResult(result);
    if (parsed?.success) {
      projectiles.value = parsed.data || [];
    } else {
      message.error(parsed?.error || t('common.error') || 'Error loading projectiles');
    }
  } catch (err) {
    message.error(err.message || 'Failed to load projectiles');
  } finally {
    loading.value = false;
  }
};

// Create projectile
const handleCreateProjectile = async () => {
  if (!createFormRef.value) return;

  try {
    await createFormRef.value.validate();
    creating.value = true;

    const fn = getBinding('ProjectileCreateProjectile');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }

    const result = await fn(
      createForm.value.name,
      createForm.value.weight,
      createForm.value.bc
    );

    if (result?.success) {
      message.success(t('common.created') || 'Projectile created');
      showCreateModal.value = false;
      createForm.value = {
        name: '',
        weight: 0,
        bc: 0.5,
      };
      await loadProjectiles();
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
const openEditModal = (projectile) => {
  editingProjectile.value = projectile;
  editForm.value = {
    name: projectile.name || '',
    weight: projectile.weightGrams || 0,
    bc: projectile.bc || 0.5,
  };
  showEditModal.value = true;
};

// Save projectile
const handleSaveProjectile = async () => {
  if (!editingProjectile.value) return;

  try {
    if (editFormRef.value) {
      await editFormRef.value.validate();
    }
    saving.value = true;

    const fn = getBinding('ProjectileUpdateProjectile');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }

    const result = await fn(
      editingProjectile.value.id,
      editForm.value.name,
      editForm.value.weight,
      editForm.value.bc
    );
    if (result?.success) {
      message.success(t('common.saved') || 'Projectile updated');
      showEditModal.value = false;
      await loadProjectiles();
    } else {
      message.error(result?.error);
    }
  } catch (err) {
    message.error(err.message);
  } finally {
    saving.value = false;
  }
};

// Delete projectile
const handleDeleteProjectile = (projectile) => {
  dialog.warning({
    title: t('common.delete') || 'Delete',
    content: `${t('common.deleteConfirm') || 'Delete'} "${projectile.name}"?`,
    positiveText: t('common.delete') || 'Delete',
    negativeText: t('common.cancel') || 'Cancel',
    onPositiveClick: async () => {
      try {
        const fn = getBinding('ProjectileDeleteProjectile');
        if (!fn) {
          message.error('Backend not ready');
          return;
        }
        const result = await fn(projectile.id);
        if (result?.success) {
          message.success(t('common.deleted') || 'Deleted');
          await loadProjectiles();
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
  await loadProjectiles();
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
