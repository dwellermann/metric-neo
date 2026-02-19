<script setup>
import { NCard, NGrid, NGi } from 'naive-ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>
<template>
  <div class="dashboard">
    <n-card>
      <template #header>
        <div class="card-header">
          <span class="mdi mdi-view-dashboard"></span>
          <span class="title">{{ t('dashboard.readyToStart') }}</span>
        </div>
      </template>
      <p>{{ t('dashboard.dataReady') }}</p>
      <p class="muted">{{ t('dashboard.laterFeatures') }}</p>
    </n-card>

    <h3 class="section-title">{{ t('dashboard.quickActions') }}</h3>
    <n-grid :cols="3" x-gap="16" y-gap="16" responsive="screen">
      <n-gi>
        <n-card class="action-card" @click="$router.push('/profiles')">
          <div class="action-content">
            <span class="mdi mdi-account-circle action-icon"></span>
            <h4>{{ t('profiles.title') }}</h4>
            <p class="muted">{{ t('profiles.description') }}</p>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="action-card" @click="$router.push('/projectiles')">
          <div class="action-content">
            <span class="mdi mdi-ammunition action-icon"></span>
            <h4>{{ t('projectiles.title') }}</h4>
            <p class="muted">{{ t('projectiles.description') }}</p>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="action-card" @click="$router.push('/sessions')">
          <div class="action-content">
            <span class="mdi mdi-chart-line action-icon"></span>
            <h4>{{ t('sessions.title') }}</h4>
            <p class="muted">{{ t('sessions.description') }}</p>
          </div>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>
<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-header .mdi {
  font-size: 1.6rem;
}

.title {
  font-weight: 600;
}

.section-title {
  margin: 0;
}

.action-card {
  cursor: pointer;
}

.action-content {
  text-align: center;
}

.action-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 8px;
}

.muted {
  opacity: 0.7;
}
</style>
