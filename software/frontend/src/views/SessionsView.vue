<template>
  <div class="page">
    <n-card>
      <template #header>
        <div class="header">
          <n-button quaternary circle @click="$router.push('/')">
            <span class="mdi mdi-arrow-left"></span>
          </n-button>
          <span class="mdi mdi-chart-line"></span>
          <span>{{ t('sessions.title') }}</span>
        </div>
      </template>
      <n-empty :description="t('sessions.title') + ' - coming soon'" />
    </n-card>
  </div>
</template>

<script setup>
import { NButton, NCard, NEmpty } from 'naive-ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<style scoped>
.page {
  padding: 16px;
}

.header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.header .mdi {
  font-size: 1.2rem;
}
</style>
