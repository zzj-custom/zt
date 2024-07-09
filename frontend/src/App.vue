<template>
  <RouterView></RouterView>

  <!-- 全局消息组件 -->
  <v-snackbar v-model="state.show" location="top" :timeout="state.timeout" :color="snackbarColor" @hide="onHide">
    {{ state.currentMessage?.content }}
    <template #actions>
      <v-btn color="white" variant="text" @click="state.show = false">Close</v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import { inject, computed } from 'vue'
import { MessageState } from './plugins/messagePlugin'

const state = inject<MessageState>('messageState')!
const displayNextMessage = inject<() => void>('displayNextMessage')!

const snackbarColor = computed(() => {
  console.log("state:",state)
  switch (state.currentMessage?.type) {
    case 'success':
      return 'green'
    case 'warning':
      return 'orange'
    case 'error':
      return 'red'
    default:
      return 'info'
  }
})

const onHide = () => {
  state.show = false
  setTimeout(displayNextMessage, 300) // Delay to show the next message
}
</script>