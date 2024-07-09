<script setup lang="ts">
import {RouterView} from 'vue-router';
import Footer from './Footer.vue';
import Header from "@/layout/Header.vue";

import {useAppStore} from "@/pinia/useAppStoe";
import Navigation from "@/layout/Navigation.vue";

const appStore = useAppStore();
</script>
<template>
  <v-layout ref="app">
    <v-navigation-drawer left elevation="0" app class="leftSidebar" v-model="appStore.drawer">
      <Navigation/>
    </v-navigation-drawer>

    <v-app-bar ref="appBar" elevation="0" height="70">
      <Header/>
    </v-app-bar>

    <v-main>
      <v-container fluid class="page-wrapper">
        <div class="maxWidth">
          <RouterView/>
        </div>
      </v-container>
    </v-main>

    <v-footer app>
      <Footer/>
    </v-footer>
  </v-layout>
</template>


<style lang="scss" scoped>
.footer {
  padding-bottom: 0 !important;
}
</style>