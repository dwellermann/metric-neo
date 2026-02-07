<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { NCard, NSelect, NSpace, NText } from 'naive-ui';

const { t, locale } = useI18n();

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
</script>

<template>
  <div class="settings-page">
    <n-card :title="t('settings.title')" :segmented="true">
      <n-space vertical>
        <div class="setting-item">
          <n-text strong>{{ t('settings.language') }}</n-text>
          <n-select
            v-model:value="currentLanguage"
            :options="languageOptions"
            style="width: 200px"
          />
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
