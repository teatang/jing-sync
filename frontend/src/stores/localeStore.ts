import { defineStore } from 'pinia'
import { type Language } from '@/types/index'
import i18n from '@/i18n'

interface LocaleStates {
    currentLang: string;
}

interface LocaleActions {
    setLanguage: (lang: Language) => void;
}

export const useLocaleStore = defineStore<string, LocaleStates, {}, LocaleActions>('locale', {
    state: () => ({
        currentLang: localStorage.getItem('lang') || 'zh-CN'
    }),
    actions: {
        setLanguage(lang: Language) {
            this.currentLang = lang
            localStorage.setItem('lang', lang)
            i18n.global.locale.value = lang
        }
    }
})
