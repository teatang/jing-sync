import { defineStore } from "pinia";
import type { IConfig, Engine, InfoList, Pagination, LoginInfo } from "@/types";
import { setToken } from "@/utils/token";
import RestClient from "@/utils/rest-client";

export const useMainStore = defineStore<string, IConfig>("main", {
  state: (): IConfig => ({
    loading: false,
  }),
  actions: {
    // define your actions
  },
});

interface AllEngineListActions {
  fetchAllEngines: () => Promise<void>;
}

export const useAllEngineListStore = defineStore<string, InfoList<Engine>, {}, AllEngineListActions>(
  "allEngineListStore",
  {
    state: () => ({
      list: [] as Engine[],
      pagination: {} as Pagination,
    }),
    actions: {
      async fetchAllEngines() {
        try {
          const infoList = await new RestClient().get<InfoList<Engine>>(
            "/engine",
            {
              page: 1,
              size: 100,
            }
          );
          this.list = infoList.list;
          this.pagination = infoList.pagination;
        } catch (error) {
          console.error("获取engines失败:", error);
        }
      },
    },
  }
);

interface UserActions {
  login: (loginForm: { username: string; password: string }) => Promise<void>;
}

export const useUserStore = defineStore<string, LoginInfo, {}, UserActions>('user', {
  state: () => ({
    token: '',
    expire_at: ''
  }),
  
  actions: {
    async login(loginForm: { username: string; password: string }) {
      try {
        const res = await new RestClient().post<LoginInfo>("/login", loginForm)
        this.token = res.token
        setToken(res.token)
      } catch (error) {
        return Promise.reject(error)
      }
    }
  }
})
