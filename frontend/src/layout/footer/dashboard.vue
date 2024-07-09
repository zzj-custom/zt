<template>
  <v-footer
      class="bg-blue_light text-center d-flex flex-column"
  >
    <div class="pt-0 text-body-1">
      该网站主要是一些小工具集合，结合自己的爱好写的一些小东西，如果对你有帮助可以使用，存在任何问题可以给我留言<v-btn size="small" icon="mdi-email" variant="text"
    ></v-btn>
    </div>
    <v-divider></v-divider>

    <div class="text-body-2">
      ©{{ new Date().getFullYear() }} All Rights Reserved. 归zk网站所有。
    </div>
  </v-footer>
</template>

<script>
export default {
  data: () => ({
    icons: [
      'mdi-facebook',
      'mdi-twitter',
      'mdi-linkedin',
      'mdi-instagram',
    ],
  }),
}
</script>