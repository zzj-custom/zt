<script setup lang="ts">
import {reactive, ref} from 'vue'
import {helpers, required} from "@vuelidate/validators";
import {useVuelidate} from "@vuelidate/core";
import {List} from "../../../wailsjs/go/app/app";
import {VideoList} from "@/types/video"

interface Variable {
  loading: boolean
  url: string
}
const variable = reactive<Variable>({
  loading: false,
  url: '',
})

const videoList = defineModel({
  type: Array as () => VideoList[],
})

const rules  = ref<any>({
  url:{
    required: helpers.withMessage('视频地址必须填写', required),
  }
})

const v$ = useVuelidate(rules, variable)

const search = async () => {
  await v$.value.$validate();
  if (v$.value.$invalid) {
    console.error('表单验证失败');
    return;
  }

  variable.loading = true;
  try {
    const response = await List(variable.url);
    if (response.code === 0) {
      videoList.value = response.result; // 直接赋值，如果 result 已经是 VideoList[]
    } else {
      console.error(response.msg);
    }
  } catch (error:any) {
    console.error('搜索过程中发生错误:', error);
  } finally {
    variable.loading = false;
  }
};



</script>
<template>
  <v-card class="mx-auto" hover>
    <v-card-item>
      <v-card-title> 搜索地址 </v-card-title>
    </v-card-item>

    <v-card-text class="text-center">
      <div>
        <div class="my-3 d-flex justify-center">
          <v-text-field
              v-model="variable.url"
              max-width="50%"
              :loading="variable.loading"
              append-inner-icon="mdi-magnify"
              density="compact"
              label="请填写需要下载的视频地址"
              variant="solo"
              :error-messages="v$.url.$errors.map((e:any) => e.$message)"
              single-line
              @click:append-inner="search"
          ></v-text-field>
        </div>
        <div class="d-flex justify-center">
          <v-chip density="compact" variant="text">目前允许下载视频包括：</v-chip>
          <v-chip
              density="compact"
              color="cyan"
              label
          >
            <v-icon icon="mdi-twitter" start></v-icon>
            New Tweets
          </v-chip>
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>