import { createI18n } from 'vue-i18n'
import en from '@/locales/en'
import zh from '@/locales/zh-CN'

declare module "vue-i18n" {
    export interface DefineLocaleMessage {
        msg: {
            delete_confirm_msg_title: string,
            delete_confirm_msg: string,
            delete_confirm_msg_success: string,
        },
        page: {
            handle: string;
            handle_edit: string;
            handle_delete: string;
            handle_add: string;
            create_time: string;
            update_time: string;
            user: {
                username: string;
                password: string;
            };
        };
    }
}

const i18n = createI18n({
    legacy: false,
    locale: localStorage.getItem('lang') || 'zh-CN',
    fallbackLocale: 'zh-CN',
    messages: {
        'zh-CN': zh,
        'en': en,
    }
})

export default i18n