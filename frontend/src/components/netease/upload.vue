<script setup lang="ts">
import { response } from "../../../wailsjs/go/models";
import { ChooseFile, ChooseFolder } from "../../../wailsjs/go/app/App";
import { NeList } from "@/types/netease";
const fileList = defineModel({
  type: Array as () => NeList[],
});

const processFiles = async () => {
  try {
    const response: response.Reply = await ChooseFile();
    if (response.code === 0) {
      fileList.value = [
        {
          ...response.result,
        },
      ];
    }
  } catch (error) {
    console.error("Error uploading files:", error);
  }
};


const processFolder = async () => {
  try {
    const response: response.Reply = await ChooseFolder();
    if (response.code === 0) {
      fileList.value = [...response.result];
    }
  } catch (error) {
    console.error("Error uploading files:", error);
  }
};
</script>

<script lang="ts">
export default {
  name: "NeUpload",
};
</script>

<template>
  <v-card class="mx-auto" hover>
    <v-card-item>
      <v-card-title> 文件上传 </v-card-title>

      <v-card-subtitle> 选择你需要转换的网易云文件 </v-card-subtitle>
    </v-card-item>

    <v-card-text class="text-center">
      <div>
        <!-- <v-row justify="center">
          <v-col cols="auto">
            <v-img
              src="https://cdn.vuetifyjs.com/images/parallax/material.jpg"
              alt="Profile Image"
              width="120"
              height="120"
              cover
              class="rounded-circle"
            ></v-img>
          </v-col>
        </v-row> -->

        <div class="d-flex justify-center my-3 gap-3">
          <v-btn class="bg-primary" @click="processFiles">选择文件</v-btn>
          <v-btn class="bg-primary" @click="processFolder">选择目录</v-btn>
          <v-hover v-slot="{ isHovering, props }"
            ><v-btn
              v-bind="props"
              :class="{ 'on-hover': isHovering, 'bg-error': isHovering }"
              class="text-none"
              :color="isHovering ? 'undeinfed' : 'error'"
              variant="outlined"
              @click="
                () => {
                  fileList = [];
                }
              "
              >重置文件</v-btn
            >
          </v-hover>
        </div>
        <p class="mb-0">Allowed NCM File. Max size of 30M</p>
      </div>
    </v-card-text>
  </v-card>
</template>
