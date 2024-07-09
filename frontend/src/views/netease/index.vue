<script setup lang="ts">
import { ref, watch,inject } from "vue";
import NeListComponent from "@/components/netease/list.vue";
import NeUpload from "@/components/netease/upload.vue";
import { NeList } from "@/types/netease";
import NeFooter from "@/components/netease/footer.vue";
import {useAppStore} from "@/pinia/useAppStoe";
import {app, response} from "../../../wailsjs/go/models";
import {Process} from "../../../wailsjs/go/app/App";

const fileList = ref<NeList[]>([]);
const selectedItem = ref<number>(1);
const outPath = ref<string>("");
const proc = ref<number>(1);
const appStore = useAppStore();

watch(outPath, (newOutPath:string) => {
  if (newOutPath !== outPath.value) {
    outPath.value = newOutPath;
  }
})

watch(fileList, (newFileList:NeList[]) => {
  if (newFileList !== fileList.value) {
    fileList.value = [...newFileList];
  }
})


watch(proc, (newProc:number) => {
  if (newProc === 2) {
    const value: app.ProcessRequest[] = [];
    for (const item of fileList.value) {
      value.push({
        flag: item.flag,
        outPath: outPath.value,
        pType:selectedItem.value,
      });
    }
    process(value);
  }
});

const process = async (value: app.ProcessRequest[]) => {
  try {
    const response: response.Reply = await Process(value);
    if (response.code === 0) {
      fileList.value = fileList.value.map((item) => {
        return {
          ...item,
          status:2
        }
      });
    }
  } catch (error) {
    console.error("Error uploading files:", error);
  }
}

</script>

<template>
  <NeUpload v-model="fileList"></NeUpload>
  <NeListComponent v-model="fileList" class="mt-3"></NeListComponent>
  <NeFooter
    elevation="0"
    :class="appStore.getDrawer() ? 'fixed-footer' : 'tf'"
    v-model:selectedItem="selectedItem"
    v-model:proc="proc"
    v-model:out-path="outPath"
  ></NeFooter>
</template>

<style sass>
.fixed-footer {
  position: fixed;
  bottom: 0;
  left: 256px;
  width: calc(100% - 256px);
  max-width: calc(
    100% - 256px
  ); /* Ensure it doesn't exceed the container width */
  background-color: white; /* Ensure the footer background matches the content */
  z-index: 999;
}

.tf{
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  max-width: 100%; /* Ensure it doesn't exceed the container width */
  background-color: white; /* Ensure the footer background matches the content */
  z-index: 999;
}
</style>
