// src/stores/auth.ts
import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null as string | null,
    user: null as Record<string, any> | null,
  }),
  actions: {
    setAuth(token: string, user: Record<string, any>) {
      this.token = token;
      this.user = user;
    },
    logout() {
      this.token = null;
      this.user = null;
    }
  },
  persist: {
    key: 'auth',
    storage: localStorage,
  },
});
