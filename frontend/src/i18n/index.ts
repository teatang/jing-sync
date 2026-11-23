import { createI18n } from 'vue-i18n'
import en from '@/locales/en'
import zh from '@/locales/zh-CN'

declare module "vue-i18n" {
    export interface DefineLocaleMessage {
        theme: {
            light: string;
            dark: string;
            auto: string;
        };
        nav: {
            jobs: string;
            engines: string;
            logs: string;
            users: string;
            logout: string;
        };
        buttons: {
            add: string;
            edit: string;
            delete: string;
            save: string;
            cancel: string;
            test: string;
            syncNow: string;
            enable: string;
            disable: string;
            refresh: string;
            details: string;
        };
        login: {
            title: string;
            subtitle: string;
            badge: string;
            username: string;
            username_placeholder: string;
            password: string;
            password_placeholder: string;
            username_required: string;
            password_required: string;
            submit: string;
            submitting: string;
            success: string;
            feature_sync: string;
            feature_jwt: string;
            feature_cron: string;
            copyright: string;
        };
        msg: {
            delete_confirm_msg_title: string;
            delete_confirm_msg: string;
            delete_confirm_msg_success: string;
            test_success: string;
            test_failed: string;
            sync_triggered: string;
            sync_failed: string;
        };
        sync: {
            mode: string;
            mode_add_only: string;
            mode_full: string;
            mode_move: string;
            speed: string;
            speed_normal: string;
            speed_fast: string;
            speed_slow: string;
            min_size: string;
            max_size: string;
            excludes: string;
            excludes_placeholder: string;
        };
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