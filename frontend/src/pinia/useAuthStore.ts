import { defineStore } from "pinia";
import { UserInfo } from "@/types/auth";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    authorization: false,
    user: {} as UserInfo,
  }),
  actions: {
      setUser(user: UserInfo) {
          this.user = user;
      },
      logout() {
          this.authorization = false;
          this.user = {Id: "", avatar: "", email: "", mobile: "", name: ""};
      },
      login(user: UserInfo) {
          this.authorization = true;
          this.setUser(user);
      }
  }
})