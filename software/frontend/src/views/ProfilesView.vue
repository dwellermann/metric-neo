<template>
  <div class="page">
    <n-card :segmented="true">
      <template #header>
        <div class="header">
          <span class="mdi mdi-pistol"></span>
          <span>{{ t('profiles.title') }}</span>
        </div>
      </template>

      <template #header-extra>
        <n-button type="primary" @click="showCreateModal = true" :disabled="loading">
          {{ t('profiles.createNew') || 'New Profile' }}
        </n-button>
      </template>

      <!-- Profiles Table -->
      <n-data-table :columns="columns" :data="profiles" :pagination="false"></n-data-table>
    </n-card>

    <!-- Create Profile Modal -->
    <n-modal v-model:show="showCreateModal" preset="dialog" :title="t('profiles.createNew') || 'New Profile'">
      <n-form ref="createFormRef" :model="createForm" :rules="formRules">
        <n-form-item :label="t('profiles.name') || 'Name'" path="name">
          <n-input v-model:value="createForm.name" placeholder="e.g., My Rifle" />
        </n-form-item>
        <n-form-item :label="t('profiles.category') || 'Category'" path="category">
          <n-select v-model:value="createForm.category" :options="categoryOptions" />
        </n-form-item>
        <n-form-item :label="t('profiles.barrelLength') || 'Barrel Length (mm)'" path="barrelLength">
          <n-input-number v-model:value="createForm.barrelLength" :min="0" />
        </n-form-item>
        <n-form-item :label="t('profiles.triggerWeight') || 'Trigger Weight (g)'" path="triggerWeight">
          <n-input-number v-model:value="createForm.triggerWeight" :min="0" />
        </n-form-item>
        <n-form-item :label="t('profiles.sightHeight') || 'Sight Height (mm)'" path="sightHeight">
          <n-input-number v-model:value="createForm.sightHeight" :min="0" />
        </n-form-item>
        <n-divider />
        <n-form-item :label="t('profiles.hasOptics') || 'Has Optics'" path="hasOptics">
          <n-switch v-model:value="createForm.hasOptics" />
        </n-form-item>
        <div v-if="createForm.hasOptics">
          <n-form-item :label="t('profiles.opticType') || 'Optic Type'" path="opticType">
            <n-select v-model:value="createForm.opticType" :options="opticTypeOptions" />
          </n-form-item>
          <n-form-item :label="t('profiles.modelName') || 'Model Name'" path="modelName">
            <n-input v-model:value="createForm.modelName" />
          </n-form-item>
          <n-form-item :label="t('profiles.opticWeight') || 'Weight (g)'" path="opticWeight">
            <n-input-number v-model:value="createForm.opticWeight" :min="0" />
          </n-form-item>
          <n-form-item :label="t('profiles.minMagnification') || 'Min Magnification'" path="minMagnification">
            <n-input-number v-model:value="createForm.minMagnification" :min="0" :step="0.1" />
          </n-form-item>
          <n-form-item :label="t('profiles.maxMagnification') || 'Max Magnification'" path="maxMagnification">
            <n-input-number v-model:value="createForm.maxMagnification" :min="0" :step="0.1" />
          </n-form-item>
        </div>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showCreateModal = false">{{ t('common.cancel') || 'Cancel' }}</n-button>
          <n-button type="primary" @click="handleCreateProfile" :loading="creating">
            {{ t('common.create') || 'Create' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- Edit Profile Modal -->
    <n-modal v-model:show="showEditModal" preset="dialog" :title="t('profiles.editProfile') || 'Edit Profile'">
      <n-form v-if="editingProfile" ref="editFormRef" :model="editForm" :rules="editFormRules">
        <n-form-item :label="t('profiles.name') || 'Name'" path="name">
          <n-input v-model:value="editForm.name" />
        </n-form-item>
        <n-form-item :label="t('profiles.category') || 'Category'" path="category">
          <n-select v-model:value="editForm.category" :options="categoryOptions" />
        </n-form-item>
        <n-form-item :label="t('profiles.barrelLength') || 'Barrel Length (mm)'" path="barrelLength">
          <n-input-number v-model:value="editForm.barrelLength" :min="0" />
        </n-form-item>
        <n-form-item :label="t('profiles.triggerWeight') || 'Trigger Weight (g)'" path="triggerWeight">
          <n-input-number v-model:value="editForm.triggerWeight" :min="0" />
        </n-form-item>
        <n-form-item :label="t('profiles.sightHeight') || 'Sight Height (mm)'" path="sightHeight">
          <n-input-number v-model:value="editForm.sightHeight" :min="0" />
        </n-form-item>
        <n-divider />
        <n-form-item :label="t('profiles.hasOptics') || 'Has Optics'" path="hasOptics">
          <n-switch v-model:value="editForm.hasOptics" />
        </n-form-item>
        <div v-if="editForm.hasOptics">
          <n-form-item :label="t('profiles.selectOptic') || 'Select Optic'" path="sightId">
            <n-select v-model:value="editForm.sightId" :options="sightOptions" clearable />
          </n-form-item>
        </div>
        <n-divider />
        <n-form-item :label="t('profiles.twistRate') || 'Twist Rate (mm)'" path="twistRate">
          <n-input-number v-model:value="editForm.twistRate" :min="0" :step="0.1" clearable />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showEditModal = false">{{ t('common.cancel') || 'Cancel' }}</n-button>
          <n-button type="primary" @click="handleSaveProfile" :loading="saving">
            {{ t('common.save') || 'Save' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, h, onMounted } from 'vue';
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
  NSpin,
  NDataTable,
  NSpace,
  NSwitch,
  NDivider,
} from 'naive-ui';
import { useDialog, useMessage } from 'naive-ui';

const { t } = useI18n();
const dialog = useDialog();
const message = useMessage();

// Helper to get Wails bindings - they're injected at runtime
const getBinding = (method) => {
  if (!window.go?.main?.App) {
    console.error('Wails bindings not available yet');
    return null;
  }
  return window.go.main.App[method];
};

// Helper to parse Wails Result data (handles double-stringification)
const parseWailsResult = (result) => {
  if (!result?.data) return result;

  // Check if data is double-stringified
  if (Array.isArray(result.data)) {
    const parsed = {
      ...result,
      data: result.data.map((item, idx) => {
        // Try to detect if it's a JSON string by looking at the stringified form
        const stringified = JSON.stringify(item);
        if (stringified.includes('\\\"') || (stringified.startsWith('"{') && stringified.endsWith('}"'))) {
          try {
            // Parse once to get the string content
            const firstParse = typeof item === 'string' ? item : JSON.parse(item);
            // Parse again to get the object
            const secondParse = JSON.parse(firstParse);
            return secondParse;
          } catch (e) {
            return item;
          }
        }

        // Single parse attempt
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

const profiles = ref([]);
const sights = ref([]);
const loading = ref(true);
const creating = ref(false);
const saving = ref(false);

const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingProfile = ref(null);

const createFormRef = ref(null);
const editFormRef = ref(null);

const categoryOptions = [
  { label: 'Air Rifle', value: 'air_rifle' },
  { label: 'Air Pistol', value: 'air_pistol' },
  { label: 'Bow', value: 'bow' },
  { label: 'Firearm', value: 'firearm' },
];

const opticTypeOptions = [
  { label: 'Scope', value: 'scope' },
  { label: 'Red Dot', value: 'red_dot' },
  { label: 'Diopter', value: 'diopter' },
  { label: 'Open Sights', value: 'open_sights' },
];

const sightOptions = computed(() =>
  sights.value.map((sight) => ({
    label: `${sight.modelName} (${sight.type})`,
    value: sight.id,
  }))
);

const createForm = ref({
  name: '',
  category: 'air_rifle',
  barrelLength: 0,
  triggerWeight: 0,
  sightHeight: 0,
  hasOptics: false,
  opticType: 'scope',
  modelName: '',
  opticWeight: 0,
  minMagnification: 1,
  maxMagnification: 10,
});

const editForm = ref({
  name: '',
  category: 'air_rifle',
  barrelLength: 0,
  triggerWeight: 0,
  sightHeight: 0,
  hasOptics: false,
  sightId: null,
  twistRate: null,
});

const formRules = {
  name: {
    required: true,
    message: t('validation.required') || 'Name is required',
    trigger: 'blur',
  },
  modelName: {
    required: (rule, value) => {
      if (createForm.value.hasOptics && !value) {
        return new Error(t('validation.required') || 'Model name is required');
      }
      return true;
    },
    trigger: 'blur',
  },
};

const editFormRules = {
  name: {
    required: true,
    message: t('validation.required') || 'Name is required',
    trigger: 'blur',
  },
  category: {
    required: true,
    message: t('validation.required') || 'Category is required',
    trigger: 'change',
  },
  sightId: {
    required: (rule, value) => {
      if (editForm.value.hasOptics && !value) {
        return new Error(t('validation.required') || 'Optic is required');
      }
      return true;
    },
    trigger: 'change',
  },
};

const columns = [
  {
    title: t('profiles.name') || 'Name',
    key: 'name',
  },
  {
    title: t('profiles.category') || 'Category',
    key: 'category',
    render: (row) => row.category?.replace(/_/g, ' ') || '-',
  },
  {
    title: t('profiles.barrelLength') || 'Barrel',
    key: 'barrelLengthMM',
    render: (row) => row.barrelLengthMM ? `${row.barrelLengthMM.toFixed(1)} mm` : '-',
  },
  {
    title: t('profiles.triggerWeight') || 'Trigger',
    key: 'triggerWeightG',
    render: (row) => row.triggerWeightG ? `${row.triggerWeightG.toFixed(1)} g` : '-',
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
            onClick: () => handleDeleteProfile(row),
          },
          { default: () => h('span', { class: 'mdi mdi-delete' }) }
        ),
      ],
    }),
  },
];

// Load profiles
const loadProfiles = async () => {
  loading.value = true;
  try {
    const fn = getBinding('ProfileListProfiles');
    if (!fn) {
      profiles.value = [];
      return;
    }
    const result = await fn();
    const parsed = parseWailsResult(result);
    if (parsed?.success) {
      profiles.value = parsed.data || [];
    } else {
      message.error(parsed?.error || t('common.error') || 'Error loading profiles');
    }
  } catch (err) {
    message.error(err.message || 'Failed to load profiles');
  } finally {
    loading.value = false;
  }
};

// Load sights
const loadSights = async () => {
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
    }
  } catch (err) {
    // ignore for now
  }
};

// Create profile
const handleCreateProfile = async () => {
  if (!createFormRef.value) return;

  try {
    await createFormRef.value.validate();
    creating.value = true;

    const fn = getBinding('ProfileCreateProfile');
    if (!fn) {
      message.error('Backend not ready');
      return;
    }

    const result = await fn(
      createForm.value.name,
      createForm.value.category,
      createForm.value.barrelLength,
      createForm.value.triggerWeight,
      createForm.value.sightHeight
    );

    if (result?.success) {
      if (createForm.value.hasOptics) {
        const setOptic = getBinding('ProfileSetOptic');
        if (!setOptic) {
          message.error('Backend not ready');
          return;
        }
        const opticResult = await setOptic(
          result.data.id,
          createForm.value.opticType,
          createForm.value.modelName,
          createForm.value.opticWeight,
          createForm.value.minMagnification,
          createForm.value.maxMagnification
        );
        if (!opticResult?.success) {
          message.error(opticResult?.error);
          return;
        }
      }
      message.success(t('common.created') || 'Profile created');
      showCreateModal.value = false;
      createForm.value = {
        name: '',
        category: 'air_rifle',
        barrelLength: 0,
        triggerWeight: 0,
        sightHeight: 0,
        hasOptics: false,
        opticType: 'scope',
        modelName: '',
        opticWeight: 0,
        minMagnification: 1,
        maxMagnification: 10,
      };
      await loadProfiles();
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
const openEditModal = (profile) => {
  editingProfile.value = profile;
  editForm.value = {
    name: profile.name || '',
    category: profile.category || 'air_rifle',
    barrelLength: profile.barrelLengthMM || 0,
    triggerWeight: profile.triggerWeightG || 0,
    sightHeight: profile.sightHeightMM || 0,
    hasOptics: !!profile.optic || !!profile.opticID,
    sightId: profile.opticID || null,
    twistRate: profile.twistRateMM || null,
  };
  showEditModal.value = true;
};

// Save profile (optics + twist rate)
const handleSaveProfile = async () => {
  if (!editingProfile.value) return;

  try {
    if (editFormRef.value) {
      await editFormRef.value.validate();
    }
    saving.value = true;

    const updateFn = getBinding('ProfileUpdateProfile');
    if (!updateFn) {
      message.error('Backend not ready');
      return;
    }

    const updateResult = await updateFn(
      editingProfile.value.id,
      editForm.value.name,
      editForm.value.category,
      editForm.value.barrelLength,
      editForm.value.triggerWeight,
      editForm.value.sightHeight
    );
    if (!updateResult?.success) {
      message.error(updateResult?.error);
      return;
    }

    // Handle optics linking/removal
    if (editForm.value.hasOptics) {
      // Benutzer will Optik, aber hat keine ausgewählt
      if (!editForm.value.sightId) {
        // Falls vorher eine Optik da war, entfernen
        if (editingProfile.value.optic || editingProfile.value.opticID) {
          const fn = getBinding('ProfileRemoveOptic');
          if (!fn) {
            message.error('Backend not ready');
            return;
          }
          const removeResult = await fn(editingProfile.value.id);
          if (!removeResult?.success) {
            message.error(removeResult?.error);
            return;
          }
        }
      } else {
        // Benutzer hat eine Optik ausgewählt
        const fn = getBinding('ProfileLinkOpticByID');
        if (!fn) {
          message.error('Backend not ready');
          return;
        }
        const opticResult = await fn(editingProfile.value.id, editForm.value.sightId);
        if (!opticResult?.success) {
          message.error(opticResult?.error);
          return;
        }
      }
    } else {
      // Benutzer hat hasOptics auf false gesetzt - Optik entfernen
      if (editingProfile.value.optic || editingProfile.value.opticID) {
        const fn = getBinding('ProfileRemoveOptic');
        if (!fn) {
          message.error('Backend not ready');
          return;
        }
        const removeResult = await fn(editingProfile.value.id);
        if (!removeResult?.success) {
          message.error(removeResult?.error);
          return;
        }
      }
    }

    // Handle twist rate
    if (editForm.value.twistRate !== null && editForm.value.twistRate !== undefined) {
      const fn = getBinding('ProfileSetTwistRate');
      if (!fn) {
        message.error('Backend not ready');
        return;
      }
      const twistResult = await fn(editingProfile.value.id, editForm.value.twistRate);
      if (!twistResult?.success) {
        message.error(twistResult?.error);
        return;
      }
    } else if (editingProfile.value.twistRateMM !== null && editingProfile.value.twistRateMM !== undefined) {
      const fn = getBinding('ProfileRemoveTwistRate');
      if (!fn) {
        message.error('Backend not ready');
        return;
      }
      const twistResult = await fn(editingProfile.value.id);
      if (!twistResult?.success) {
        message.error(twistResult?.error);
        return;
      }
    }

    message.success(t('common.saved') || 'Profile updated');
    showEditModal.value = false;
    await loadProfiles();
  } catch (err) {
    message.error(err.message);
  } finally {
    saving.value = false;
  }
};

// Delete profile
const handleDeleteProfile = (profile) => {
  dialog.warning({
    title: t('common.delete') || 'Delete',
    content: `${t('common.deleteConfirm') || 'Delete'} "${profile.name}"?`,
    positiveText: t('common.delete') || 'Delete',
    negativeText: t('common.cancel') || 'Cancel',
    onPositiveClick: async () => {
      try {
        const fn = getBinding('ProfileDeleteProfile');
        if (!fn) {
          message.error('Backend not ready');
          return;
        }
        const result = await fn(profile.id);
        if (result?.success) {
          message.success(t('common.deleted') || 'Deleted');
          await loadProfiles();
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
  await loadProfiles();
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
