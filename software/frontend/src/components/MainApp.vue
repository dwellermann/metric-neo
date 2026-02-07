<script setup>
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import {
  NButton,
  NLayout,
  NLayoutContent,
  NLayoutSider,
  NMenu,
  NDivider,
} from 'naive-ui';

const { t, locale } = useI18n();
const router = useRouter();
const route = useRoute();
const collapsed = ref(false);

const mainMenuOptions = computed(() => [
  {
    label: t('dashboard.quickActions'),
    key: '/',
    icon: () => h('span', { class: 'mdi mdi-view-dashboard' }),
  },
  {
    label: t('profiles.title'),
    key: '/profiles',
    icon: () => h('span', { class: 'mdi mdi-pistol' }),
  },
  {
    label: t('projectiles.title'),
    key: '/projectiles',
    icon: () => h('span', { class: 'mdi mdi-ammunition' }),
  },
  {
    label: t('sights.title'),
    key: '/sights',
    icon: () => h('span', { class: 'mdi mdi-crosshairs-gps' }),
  },
  {
    label: t('sessions.title'),
    key: '/sessions',
    icon: () => h('span', { class: 'mdi mdi-chart-line' }),
  },
]);

const settingsMenuOptions = computed(() => [
  {
    label: t('settings.title'),
    key: '/settings',
    icon: () => h('span', { class: 'mdi mdi-cog' }),
  },
]);

function handleSelect(key) {
  if (key !== route.path) {
    router.push(key);
  }
}
</script>

<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed="collapsed"
      :collapsed-width="64"
      :width="240"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
    >
      <div style="display: flex; flex-direction: column; height: 100%">
        <!-- Header -->
        <div style="padding: 16px; border-bottom: 1px solid var(--n-border-color); flex-shrink: 0">
          <div :style="{
            display: 'flex',
            alignItems: 'center',
            gap: '12px',
            justifyContent: collapsed ? 'center' : 'flex-start'
          }">
            <img src="../assets/logo.svg" alt="Metric Neo Logo" style="width: 36px; height: 36px; flex-shrink: 0" />
            <span v-if="!collapsed" style="font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis">
              {{ t('app.title') }}
            </span>
          </div>
        </div>

        <!-- Main Menu (scrollable, takes remaining space) -->
        <div style="flex: 1; overflow-y: auto; min-height: 0">
          <n-menu
            :value="route.path"
            :options="mainMenuOptions"
            :collapsed="collapsed"
            :collapsed-width="64"
            :collapsed-icon-size="22"
            @update:value="handleSelect"
          />
        </div>

        <!-- Settings Menu (fixed at bottom) -->
        <div style="flex-shrink: 0; border-top: 1px solid var(--n-border-color)">
          <n-menu
            :value="route.path"
            :options="settingsMenuOptions"
            :collapsed="collapsed"
            :collapsed-width="64"
            :collapsed-icon-size="22"
            @update:value="handleSelect"
          />
        </div>
      </div>
    </n-layout-sider>

    <n-layout>
      <n-layout-content style="padding: 16px">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
/* Minimal CSS - most layout handled by Naive UI components */
</style>
