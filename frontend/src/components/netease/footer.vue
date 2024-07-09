<script setup lang="ts">
import { watch, ref } from "vue";

import { ChooseDirectory } from "../../../wailsjs/go/app/App";
import { SelectFolderList } from "@/types/netease";

const items: SelectFolderList[] = [
  {
    title: "与源文件夹相同",
    value: 1,
  },
  {
    title: "自定义文件夹",
    value: 2,
  },
];

const selectedItem = defineModel("selectedItem", { type: Number });
const proc = defineModel("proc", { type: Number, default: 1 });
const outPath = defineModel("outPath", { type: String });
const show = ref<boolean>(false);

watch(selectedItem, async (newSelectedItem) => {
  if (newSelectedItem === 2) {
    try {
      const res = await ChooseDirectory();
      outPath.value = res.result;
    } catch (err: any) {
      console.log(err);
    }
  }
});

const process = () => {
  if (proc.value === 2) {
    show.value = true;
  } else {
    proc.value = 2;
  }
};

const close = () => {
  show.value = false;
  proc.value = 1;
}
</script>
<script lang="ts">
export default {
  name: "NeFooter",
};
</script>

<template>
  <v-card>
    <v-card-text class="text-center">
      <div class="px-3">
        <v-row justify="space-between">
          <v-col cols="3">
            <v-select
              v-model="selectedItem"
              :items="items"
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
            <v-btn @click="process" class="bg-error" variant="flat"
              >{{ proc === 1 ? "全部提取" : "结束任务" }}
            </v-btn>
            <v-dialog max-width="300" v-model="show">
              <v-card title="注意">
                <v-card-text class="text-center">
                  <p>当前正在处理中！</p>
                  <p>请在处理完成后再进行此操作。</p>
                </v-card-text>
                <v-card-actions>
                  <v-btn
                    text="确定"
                    class="bg-primary"
                    @click="close"
                  ></v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-col>
        </v-row>
      </div>
    </v-card-text>
  </v-card>
</template>
