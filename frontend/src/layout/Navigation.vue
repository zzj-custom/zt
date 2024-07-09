<script setup lang="ts">

import { ref, shallowRef,provide } from 'vue';
import sidebarItems from './vertical-sidebar/sidebarItem';
import NavGroup from './vertical-sidebar/NavGroup/index.vue';
import NavItem from './vertical-sidebar/NavItem/index.vue';
import ExtraBox from './vertical-sidebar/extrabox/ExtraBox.vue';
import Logo from './logo/Logo.vue';
const sidebarMenu = shallowRef(sidebarItems);
</script>

<template>
  <!------Sidebar-------->
  <!---Logo part -->
  <div class="pa-5">
    <Logo />
  </div>
  <!-- ---------------------------------------------- -->
  <!---Navigation -->
  <!-- ---------------------------------------------- -->
  <div>
    <perfect-scrollbar class="scrollnavbar">
      <v-list class="pa-6">
        <!---Menu Loop -->
        <template v-for="(item, i) in sidebarMenu">
          <!---Item Sub Header -->
          <NavGroup :item="item" v-if="item.header" :key="item.title" />

          <!---Single Item-->
          <NavItem :item="item" v-else class="leftPadding" />
          <!---End Single Item-->
        </template>
      </v-list>
      <div class="pa-4">
        <ExtraBox />
      </div>
    </perfect-scrollbar>
  </div>
</template>

<style scoped>

</style>