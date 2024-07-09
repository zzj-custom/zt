<script setup lang="ts">
import { ref } from 'vue';
const active = ref(false);
</script>

<template>
  <v-card>
    <v-card-text class="text-center">
      <v-row justify="space-between">
        <v-col cols="3">
          <v-select
              density="compact"
              item-title="title"
              item-value="value"
              prepend-inner-icon="mdi-folder"
              variant="solo-filled"
              placeholder="请选择需要保存的文件夹"
              bg-color="primary"
              item-color="primary"
              hide-details
          >
          </v-select>
        </v-col>
        <v-col cols="2">
          <v-btn class="bg-error" variant="flat"
          >下载
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<!--<style scoped>-->
<!--.v-card-text{-->
<!--  padding: 0px !important;-->
<!--}-->
<!--</style>-->