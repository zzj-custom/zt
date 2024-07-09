<script setup lang="ts">
import NotificationDD from "@/layout/vertical-header/NotificationDD.vue";
import {Menu2Icon} from "vue-tabler-icons";
import ProfileDD from "@/layout/vertical-header/ProfileDD.vue";
import {useAppStore}   from "@/pinia/useAppStoe";
const appStore = useAppStore();
</script>

<template>
  <div class="d-flex align-center justify-space-between w-100">
    <div>
      <v-btn class="hidden-lg-and-up ms-md-3 ms-sm-5 ms-3 text-muted" @click="appStore.drawer = !appStore.drawer" icon variant="flat"
             size="small">
        <Menu2Icon size="20" stroke-width="1.5" />
      </v-btn>
      <NotificationDD />
    </div>
    <div>
      <ProfileDD />
    </div>
  </div>
</template>

<style scoped>

</style>