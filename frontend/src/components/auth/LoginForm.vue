<script setup lang="ts">
import {ref, computed, reactive} from 'vue';
import {LoginForm} from '@/types/auth'
import {Captcha} from '../../../wailsjs/go/app/App'
import { useVuelidate } from '@vuelidate/core'
import { required, minLength,helpers } from '@vuelidate/validators';
import {Login} from '../../../wailsjs/go/app/App'
import {app} from "../../../wailsjs/go/models";
import {useAuthStore} from "@/pinia/useAuthStore";
import router from "@/router";
const authStore = useAuthStore();

const form = reactive<LoginForm>({} as LoginForm);

const emailIdentify = (value:string) => {
  const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
  return pattern.test(value)
}

const captchaIdentify = (value:number) => {
  const pattern = /^\d{6}$/
  return pattern.test(value.toString())
}

const rules = ref({
  email: {
    required:helpers.withMessage("邮箱必须填写", required),
    emailIdentify:helpers.withMessage("请填写正确的邮箱地址", emailIdentify)
  },
  captcha: {
    required: helpers.withMessage("验证码必须填写", required),
    captchaIdentify: helpers.withMessage("请填写正确的验证码", captchaIdentify),
    minLength: helpers.withMessage("验证码长度必须为6位", minLength(6))
  }
});


// 使用 ref 创建响应式状态
const value = ref('');

// 使用 ref 创建自定义状态，初始值为 false
const custom = ref(false);

// 使用 computed 创建计算属性
const progress = computed(() => Math.min(100, value.value.length * 10));

// 使用 computed 创建颜色计算属性
const color = computed(() => {
  const baseColorIndex = Math.floor(progress.value / 40);
  const colorMap = ['error', 'warning', 'success'];
  return colorMap[baseColorIndex];
});

const captcha = async () => {
  custom.value = true;
  if (form.email === '') {
    //TODO 处理错误
    return
  }
  const response = await Captcha(form.email);
  console.log("response:",response);
  if (response.code === 0) {
    custom.value = false;
  }
}

const v$ = useVuelidate(rules, form);

const submit = async () => {
  await v$.value.$validate();
  if (!v$.value.$invalid) {
    console.log("验证成功");
    await login();
    return;
  }
  console.error("验证错误：",v$.value.$errors);
}

const login = async () => {
  const params: app.LoginRequest = {
    email:form.email,
    captcha:Number(form.captcha),
  }
  const {code, msg,result} = await Login(params);
  console.log("code:", code, "result", result, "msg", msg);
  if (code === 0) {
    authStore.login({
      Id:result.Id,
      name:result.Name,
      email:result.email,
      avatar:result.avatar ?? "@/assets/images/users/avatar-1.jpg",
      mobile:result.mobile,
    });
    await router.push("/");
  }
}

</script>

<template>
  <v-row class="d-flex mb-3" v-model="form">
    <v-col cols="12">
      <v-label class="font-weight-bold mb-1">账号</v-label>
      <v-text-field
          v-model="form.email"
          color="primary"
          density="compact"
          hite="请填写邮箱地址"
          prepend-inner-icon="mdi-email-outline"
          variant="outlined"
          :error-messages="v$.email.$errors.map(e => String(e.$message))"
          @blur="v$.email.$touch"
          @input="v$.email.$touch"
          required
          clearable
      ></v-text-field>
    </v-col>
    <v-col cols="12">
      <v-label class="font-weight-bold mb-1">密码</v-label>
      <v-text-field
          v-model="form.captcha"
          color="primary"
          density="compact"
          hint="请输入密码"
          prepend-inner-icon="mdi-lock-outline"
          variant="outlined"
          :error-messages="v$.captcha.$errors.map(e => String(e.$message))"
          @blur="v$.captcha.$touch"
          @input="v$.captcha.$touch"
          :disabled="custom"
          required
      >
      </v-text-field>
    </v-col>
    <v-col cols="12" class="pt-0">
      <div class="d-flex flex-wrap align-center ml-n2">
        <div class="ml-sm-auto">
          <RouterLink to="/"
                      class="text-primary text-decoration-none text-body-1 opacity-1 font-weight-medium">Forgot
            Password ?</RouterLink>
        </div>
      </div>
    </v-col>
    <v-col cols="12" class="pt-0">
      <v-btn @click="submit" color="primary" size="large" block type="submit">登录</v-btn>
    </v-col>
  </v-row>
</template>
