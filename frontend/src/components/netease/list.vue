<script setup lang="ts">
import { computed, ref } from "vue";
import type { NeList } from "@/types/netease";
import {Process, View} from "../../../wailsjs/go/app/App";
import {app, response} from "../../../wailsjs/go/models";
const search = ref("");

const fileList = defineModel({
  type: Array as () => NeList[],
  default: [],
});

const headers = ref<Array<Object>>([
  { title: "文件名称", key: "name" },
  { title: "音乐图片", key: "albumPic" },
  { title: "音乐名称", key: "musicName" },
  { title: "音乐ID", key: "musicId" },
  { title: "时长", key: "duration" },
  { title: "操作", key: "format", width: "100px" },
]);

const duration = computed(() => {
  {
    return (seconds: number) => {
      seconds = Math.floor(seconds / 1000); // 因为返回的是毫秒
      const hours = Math.floor(seconds / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);
      const secs = seconds % 60;

      const formattedHours = hours.toString().padStart(2, "0");
      const formattedMinutes = minutes.toString().padStart(2, "0");
      const formattedSecs = secs.toString().padStart(2, "0");

      if (hours > 0) {
        return `${formattedHours}:${formattedMinutes}:${formattedSecs}`;
      }
      return `${formattedMinutes}:${formattedSecs}`;
    };
  }
});

const deleteItem = (value: NeList) => {
  const index = fileList.value.indexOf(value);
  if (index > -1) {
    const newList = [...fileList.value];
    newList.splice(index, 1);
    fileList.value = newList;
  }
};


const extract = async (value: NeList) => {
  try {
    const params : app.ProcessRequest[] = [{
      flag:value.flag,
      outPath:"",
      pType:1
    }]
    const response: response.Reply = await Process(params);
    if (response.code === 0) {
      const resp = response.result[0]
      fileList.value = fileList.value.map((item) => {
        if (item.flag === resp.flag) {
          item.status = 2;
        }
        return item;
      })

    }
  } catch (error) {
    console.error("Error uploading files:", error);
  }
}

const view = async (value: NeList) => {
  try {
    const response: response.Reply = await View(value.flag);
    if (response.code === 0) {
      console.log("查看成功")
    }
  } catch (error) {
    console.error("Error uploading files:", error);
  }
}
</script>

<script lang="ts">
export default {
  name: "NeListComponent",
};
</script>

<template>
  <v-card hover>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon color="primary" icon="mdi-file"></v-icon> &nbsp; 文件列表

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

    <v-data-table
      v-model:search="search"
      v-model:items="fileList"
      :headers="headers"
    >
      <template v-slot:item.albumPic="{ item }">
        <v-card class="my-2" elevation="2" rounded width="64">
          <v-img :src="item.albumPic" aspect-ratio="1/1"></v-img>
        </v-card>
      </template>

      <template v-slot:item.duration="{ item }">
        <p>{{ duration(item.duration) }}</p>
      </template>

      <template v-slot:item.format="{ item }">
        <div class="d-flex justify-center gap-2">
          <v-btn density="comfortable" v-if="item.status === 1" @click="extract(item)" class="bg-primary"
            >提取</v-btn
          >
          <v-btn density="comfortable" v-else @click="view(item)" class="bg-info"
          >查看</v-btn
          >
          <v-btn
            density="comfortable"
            class="bg-error"
            @click="deleteItem(item)"
            >删除</v-btn
          >
        </div>
      </template>
    </v-data-table>
  </v-card>
</template>
