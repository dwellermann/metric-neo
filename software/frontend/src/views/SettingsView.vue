<script setup>
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSelect,
  NSpace,
  NSwitch,
  NText,
  useMessage,
} from 'naive-ui';

const { t, locale } = useI18n();
const message = useMessage();

const languageOptions = computed(() => [
  {
    label: 'Deutsch',
    value: 'de',
  },
  {
    label: 'English',
    value: 'en',
  },
]);

const currentLanguage = computed({
  get: () => locale.value,
  set: (val) => {
    locale.value = val;
  },
});

const chronoForm = ref({
  enabled: false,
  port: '',
  baudRate: 9600,
  autoRecord: true,
});

const loadingChrono = ref(false);

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

const loadChronoConfig = async () => {
  const fn = getBinding('GetChronoConfig');
  if (!fn) return;
  const result = await fn();
  const parsed = parseWailsResult(result);
  if (parsed?.success && parsed.data) {
    chronoForm.value = {
      enabled: !!parsed.data.enabled,
      port: parsed.data.port || '',
      baudRate: parsed.data.baudRate || 9600,
      autoRecord: parsed.data.autoRecord !== false,
    };
  }
};

const saveChronoConfig = async () => {
  const fn = getBinding('UpdateChronoConfig');
  if (!fn) {
    message.error('Backend not ready');
    return;
  }
  try {
    loadingChrono.value = true;
    const result = await fn(
      chronoForm.value.enabled,
      chronoForm.value.port,
      chronoForm.value.baudRate,
      chronoForm.value.autoRecord
    );
    const parsed = parseWailsResult(result);
    if (parsed?.success) {
      message.success(t('common.saved') || 'Saved');
    } else {
      message.error(parsed?.error || t('common.error') || 'Error');
    }
  } catch (err) {
    message.error(err.message || 'Error');
  } finally {
    loadingChrono.value = false;
  }
};

onMounted(async () => {
  await new Promise(resolve => setTimeout(resolve, 100));
  await loadChronoConfig();
});
</script>

<template>
  <div class="settings-page">
    <n-card :title="t('settings.title')" :segmented="true">
      <n-space vertical size="large">
        <div class="setting-item">
          <n-text strong>{{ t('settings.language') }}</n-text>
          <n-select
            v-model:value="currentLanguage"
            :options="languageOptions"
            style="width: 200px"
          />
        </div>

        <div>
          <n-text strong>{{ t('settings.chronoTitle') || 'Chronograph (RS232)' }}</n-text>
          <n-form :model="chronoForm" style="margin-top: 12px;">
            <n-form-item :label="t('settings.chronoEnabled') || 'Enabled'">
              <n-switch v-model:value="chronoForm.enabled" />
            </n-form-item>
            <n-form-item :label="t('settings.chronoPort') || 'Port'">
              <n-input v-model:value="chronoForm.port" placeholder="/dev/ttyUSB0" />
            </n-form-item>
            <n-form-item :label="t('settings.chronoBaud') || 'Baud Rate'">
              <n-input-number v-model:value="chronoForm.baudRate" :min="1200" :step="1200" />
            </n-form-item>
            <n-form-item :label="t('settings.chronoAutoRecord') || 'Auto Record'">
              <n-switch v-model:value="chronoForm.autoRecord" />
            </n-form-item>
            <n-button type="primary" @click="saveChronoConfig" :loading="loadingChrono">
              {{ t('common.save') || 'Save' }}
            </n-button>
          </n-form>
        </div>
      </n-space>
    </n-card>
  </div>
</template>

<style scoped>
.settings-page {
  max-width: 600px;
}

.setting-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
}
</style>
