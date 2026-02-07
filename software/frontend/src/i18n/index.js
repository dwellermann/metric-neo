import { createI18n } from 'vue-i18n';
import de from './locales/de.json';
import en from './locales/en.json';

// Browser-Sprache oder Default
const getBrowserLocale = () => {
  const navigatorLocale = navigator.language || navigator.userLanguage;
  if (navigatorLocale.startsWith('de')) return 'de';
  return 'en';
};

const i18n = createI18n({
  legacy: false, // Composition API
  locale: getBrowserLocale(),
  fallbackLocale: 'en',
  messages: {
    de,
    en
  }
});

export default i18n;
