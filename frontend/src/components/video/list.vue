<script setup lang="ts">
import {ref, defineModel} from "vue";
import {VideoList} from "@/types/video";

import {BrandBilibiliIcon} from 'vue-tabler-icons';

const search = ref("");

const videoList = defineModel({
  type: Array as () => VideoList[],
});

const headers = ref<Array<Object>>([
  {title: "名称", key: "title"},
  {title: "网站地址", key: "site", sortable: false},
  {title: "类型", key: "type",sortable: false},
]);

const streamsHeaders = ref<Array<Object>>([
    {title: "清晰度", key: "quality"},
    {title: "大小", key: "size"},
    {title: "格式", key: "ext", sortable: false},
]);

const selected = ref<Array<VideoList>>([]);

</script>
<template>
  <v-card class="overflow-auto" max-height="500px">
    <v-card-title class="d-flex align-center pe-2">
      <v-icon color="primary" icon="mdi-video"></v-icon> &nbsp;视频列表
      <v-spacer></v-spacer>
      <v-text-field
          v-model="search"
          density="compact"
          label="Search"
          prepend-inner-icon="mdi-magnify"
          variant="solo-filled"
          flat
          hide-details
          single-line
      ></v-text-field>
    </v-card-title>
    <v-divider></v-divider>

    <v-data-table-virtual
        v-model:search="search"
        v-model:selected="selected"
        :items="videoList"
        :headers="headers"
        item-value="name"
        show-expand
        show-select
    >
      <template v-slot:item.title="{item}">
        <v-tooltip :text="item.title">
          <template v-slot:activator="{ props }">
            <span class="text-truncate overflow-hidden" v-bind="props">{{ item.title }}</span>
          </template>
        </v-tooltip>
      </template>
      <template v-slot:header.data-table-select="{ allSelected, selectAll, someSelected }">
        <v-checkbox-btn
            :indeterminate="someSelected && !allSelected"
            :model-value="allSelected"
            color="primary"
            @update:model-value="selectAll(!allSelected)"
        ></v-checkbox-btn>
      </template>

      <template v-slot:item.data-table-select="{ internalItem, isSelected, toggleSelect }">
        <v-checkbox-btn
            :model-value="isSelected(internalItem)"
            color="primary"
            @update:model-value="toggleSelect(internalItem)"
        ></v-checkbox-btn>
      </template>
      <template v-slot:item.site="{item}">
        <v-chip
            density="compact"
            color="cyan"
            target="_blank"
            :href="item.url"
            label
        >
          <BrandBilibiliIcon class="icon"/>
          {{item.site}}
        </v-chip>
      </template>

      <template v-slot:expanded-row="{item}">
        <td  :colspan="headers.length-1">
          <v-data-table
              :headers="streamsHeaders"
              :items="item.streams"
              select-strategy="single"
              hide-default-footer
              show-select
          >

            <template v-slot:item.data-table-select="{ internalItem, isSelected, toggleSelect }">
              <v-checkbox-btn
                  :model-value="isSelected(internalItem)"
                  color="primary"
                  @update:model-value="toggleSelect(internalItem)"
              ></v-checkbox-btn>
            </template>
            <template v-slot:[`item.quality`]="{ item }">
              {{ item.quality }}
            </template>
            <template v-slot:[`item.size`]="{ item }">
              {{ (item.size / 1048576).toFixed(2) }} MB
            </template>
            <template v-slot:[`item.ext`]="{ item }">
              {{ item.ext }}
            </template>
          </v-data-table>
        </td>
      </template>
    </v-data-table-virtual>
  </v-card>
</template>
<style sass>
.icon {
  width: 24px; /* 和 v-icon 的大小一致 */
  height: 24px; /* 和 v-icon 的大小一致 */
  margin-right: 8px; /* 添加与文本之间的间距 */
  vertical-align: middle; /* 确保图标和文本垂直对齐 */
}

</style>

