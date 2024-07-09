<script setup lang="ts">
import {useAuthStore} from '@/pinia/useAuthStore'
import { computed,ref } from 'vue';
const authStore = useAuthStore()
import defaultAvatar from '@/assets/images/users/avatar-1.jpg';
import router from "@/router";
const avatar = computed(() => {
  return authStore.user?.avatar || defaultAvatar;
})

const menu = ref(false)

const logout = () => {
  useAuthStore().logout()
  router.push('/auth/login')
}

</script>

<template>
    <!-- ---------------------------------------------- -->
    <!-- notifications DD -->
    <!-- ---------------------------------------------- -->
  <v-menu
      v-model="menu"
      location="top start"
      origin="top start"
      transition="scale-transition"
      open-on-hover
  >
    <template v-slot:activator="{ props }">
      <v-chip
          v-bind="props"
          color="cyan"
          link
          pill
          class="mr-3"
      >
        <v-avatar start>
          <v-img :src="avatar"></v-img>
        </v-avatar>
        <span class="text-truncate text-h6">{{authStore.user?.name || "游客"}}</span>
      </v-chip>
    </template>

    <v-card max-width="300">
      <v-list bg-color="black">
        <v-list-item>
          <template v-slot:prepend>
            <v-avatar
                :image="avatar"
            ></v-avatar>
          </template>

          <v-list-item-title>{{authStore.user?.name || "游客"}}</v-list-item-title>

          <v-list-item-subtitle>{{authStore.user?.email || "1844066417@qq.com"}}</v-list-item-subtitle>
        </v-list-item>
      </v-list>

      <v-list>
        <v-list-item color="cyan" to="/auth/login" prepend-icon="mdi-account">
          <v-list-item-subtitle>个人信息</v-list-item-subtitle>
        </v-list-item>
      </v-list>

      <v-list>
        <v-list-item color="cyan" to="/auth/login" prepend-icon="mdi-mail">
          <v-list-item-subtitle>联系我们</v-list-item-subtitle>
        </v-list-item>
      </v-list>

      <div class="pt-4 pb-4 px-2 text-center">
        <v-btn @click="logout" color="primary" variant="outlined" block>退出</v-btn>
      </div>
    </v-card>
  </v-menu>
</template>

<style sass>
.text-cyan{
  color :#00bcd4 !important
}
</style>